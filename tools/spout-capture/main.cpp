#include <windows.h>
#include <wincodec.h>

#include <chrono>
#include <algorithm>
#include <cctype>
#include <cstddef>
#include <cstdint>
#include <ctime>
#include <filesystem>
#include <iostream>
#include <cmath>
#include <memory>
#include <sstream>
#include <string>
#include <utility>
#include <thread>
#include <vector>

#include "SpoutLibrary.h"

namespace {

struct SenderInfo {
  std::string name;
  unsigned int width = 0;
  unsigned int height = 0;
  std::string host_path;
};

struct Options {
  bool list_senders = false;
  bool capture = false;
  bool show_version = false;
  std::string sender;
  std::filesystem::path output;
  int timeout_ms = 10000;
};

struct FrameStats {
  int samples = 0;
  double mean = 0.0;
  double stddev = 0.0;
  double near_white_ratio = 0.0;
  double near_black_ratio = 0.0;
};

std::string json_escape(const std::string &value) {
  std::ostringstream out;
  for (unsigned char c : value) {
    switch (c) {
    case '"':
      out << "\\\"";
      break;
    case '\\':
      out << "\\\\";
      break;
    case '\b':
      out << "\\b";
      break;
    case '\f':
      out << "\\f";
      break;
    case '\n':
      out << "\\n";
      break;
    case '\r':
      out << "\\r";
      break;
    case '\t':
      out << "\\t";
      break;
    default:
      if (c < 0x20) {
        out << "\\u";
        const char *hex = "0123456789abcdef";
        out << "00" << hex[(c >> 4) & 0x0f] << hex[c & 0x0f];
      } else {
        out << c;
      }
    }
  }
  return out.str();
}

std::string helper_version() {
#ifdef SPOUT_CAPTURE_VERSION
  return SPOUT_CAPTURE_VERSION;
#else
  return "unknown";
#endif
}

void print_help() {
  std::cout << "spout-capture --help\n"
            << "spout-capture --version\n"
            << "spout-capture --list-senders\n"
            << "spout-capture --capture [--sender name] --output file.png --timeout-ms 10000\n";
}

void print_version() {
  std::cout << "{\"ok\":true,\"name\":\"spout-capture\",\"version\":\""
            << json_escape(helper_version()) << "\"}\n";
}

void print_error(const std::string &code, const std::string &message, const std::vector<SenderInfo> &senders = {}) {
  std::cout << "{\"ok\":false,\"code\":\"" << json_escape(code)
            << "\",\"message\":\"" << json_escape(message) << "\"";
  if (!senders.empty()) {
    std::cout << ",\"senders\":[";
    for (size_t i = 0; i < senders.size(); ++i) {
      if (i > 0) {
        std::cout << ",";
      }
      std::cout << "{\"name\":\"" << json_escape(senders[i].name) << "\",\"width\":"
                << senders[i].width << ",\"height\":" << senders[i].height
                << ",\"hostPath\":\"" << json_escape(senders[i].host_path) << "\"}";
    }
    std::cout << "]";
  }
  std::cout << "}\n";
}

bool parse_int(const std::string &value, int *out) {
  try {
    size_t used = 0;
    int parsed = std::stoi(value, &used, 10);
    if (used != value.size()) {
      return false;
    }
    *out = parsed;
    return true;
  } catch (...) {
    return false;
  }
}

bool parse_args(int argc, char **argv, Options *options, std::string *error) {
  for (int i = 1; i < argc; ++i) {
    std::string arg = argv[i];
    if (arg == "--list-senders") {
      options->list_senders = true;
    } else if (arg == "--capture") {
      options->capture = true;
    } else if (arg == "--version") {
      options->show_version = true;
    } else if (arg == "--sender") {
      if (++i >= argc) {
        *error = "--sender requires a value";
        return false;
      }
      options->sender = argv[i];
    } else if (arg == "--output") {
      if (++i >= argc) {
        *error = "--output requires a value";
        return false;
      }
      options->output = std::filesystem::path(argv[i]);
    } else if (arg == "--timeout-ms") {
      if (++i >= argc || !parse_int(argv[i], &options->timeout_ms)) {
        *error = "--timeout-ms requires an integer value";
        return false;
      }
    } else if (arg == "--help" || arg == "-h") {
      print_help();
      std::exit(0);
    } else {
      *error = "unknown argument: " + arg;
      return false;
    }
  }
  if (options->show_version) {
    print_version();
    std::exit(0);
  }
  if (options->list_senders == options->capture) {
    *error = "specify exactly one of --list-senders or --capture";
    return false;
  }
  if (options->capture && options->output.empty()) {
    *error = "--capture requires --output";
    return false;
  }
  if (options->timeout_ms < 100) {
    options->timeout_ms = 100;
  }
  return true;
}

std::string wide_to_utf8(const std::wstring &value) {
  if (value.empty()) {
    return {};
  }
  int size = WideCharToMultiByte(CP_UTF8, 0, value.c_str(), -1, nullptr, 0, nullptr, nullptr);
  if (size <= 0) {
    return {};
  }
  std::string out(static_cast<size_t>(size - 1), '\0');
  WideCharToMultiByte(CP_UTF8, 0, value.c_str(), -1, out.data(), size, nullptr, nullptr);
  return out;
}

std::string lower_copy(std::string value) {
  std::transform(value.begin(), value.end(), value.begin(), [](unsigned char c) {
    return static_cast<char>(std::tolower(c));
  });
  return value;
}

FrameStats analyze_rgba_frame(const std::vector<unsigned char> &rgba, unsigned int width, unsigned int height) {
  FrameStats stats;
  if (width == 0 || height == 0 || rgba.empty()) {
    return stats;
  }
  unsigned int step_x = 1;
  unsigned int step_y = 1;
  const unsigned int max_samples = 16384;
  while ((width / step_x) * (height / step_y) > max_samples) {
    if (width / step_x >= height / step_y) {
      step_x++;
    } else {
      step_y++;
    }
  }
  double sum = 0.0;
  double sum_sq = 0.0;
  double near_white = 0.0;
  double near_black = 0.0;
  for (unsigned int y = 0; y < height; y += step_y) {
    for (unsigned int x = 0; x < width; x += step_x) {
      const size_t i = (static_cast<size_t>(y) * width + x) * 4;
      if (i + 2 >= rgba.size()) {
        continue;
      }
      const double r = static_cast<double>(rgba[i]);
      const double g = static_cast<double>(rgba[i + 1]);
      const double b = static_cast<double>(rgba[i + 2]);
      const double luma = 0.2126 * r + 0.7152 * g + 0.0722 * b;
      sum += luma;
      sum_sq += luma * luma;
      if (luma >= 250.0) {
        near_white++;
      }
      if (luma <= 5.0) {
        near_black++;
      }
      stats.samples++;
    }
  }
  if (stats.samples == 0) {
    return stats;
  }
  const double samples = static_cast<double>(stats.samples);
  stats.mean = sum / samples;
  const double variance = std::max(0.0, (sum_sq / samples) - (stats.mean * stats.mean));
  stats.stddev = std::sqrt(variance);
  stats.near_white_ratio = near_white / samples;
  stats.near_black_ratio = near_black / samples;
  return stats;
}

bool is_blank_frame(const FrameStats &stats) {
  if (stats.samples == 0) {
    return true;
  }
  if (stats.near_black_ratio > 0.99 && stats.stddev < 1.5) {
    return true;
  }
  if (stats.near_white_ratio > 0.99 && stats.stddev < 1.5) {
    return true;
  }
  return false;
}

void force_opaque_alpha(std::vector<unsigned char> &rgba) {
  for (size_t i = 3; i < rgba.size(); i += 4) {
    rgba[i] = 255;
  }
}

bool write_png_wic(const std::filesystem::path &path, unsigned int width, unsigned int height,
                   const std::vector<unsigned char> &rgba, std::string *error) {
  IWICImagingFactory *factory = nullptr;
  IWICBitmapEncoder *encoder = nullptr;
  IWICBitmapFrameEncode *frame = nullptr;
  IWICStream *stream = nullptr;
  bool com_initialized = false;
  HRESULT hr = CoInitializeEx(nullptr, COINIT_MULTITHREADED);
  if (SUCCEEDED(hr)) {
    com_initialized = true;
  } else if (hr != RPC_E_CHANGED_MODE) {
    *error = "CoInitializeEx failed";
    return false;
  }
  auto cleanup = [&]() {
    if (frame) frame->Release();
    if (encoder) encoder->Release();
    if (stream) stream->Release();
    if (factory) factory->Release();
    if (com_initialized) CoUninitialize();
  };
  hr = CoCreateInstance(CLSID_WICImagingFactory, nullptr, CLSCTX_INPROC_SERVER,
                        IID_PPV_ARGS(&factory));
  if (FAILED(hr)) {
    cleanup();
    *error = "WIC factory creation failed";
    return false;
  }
  hr = factory->CreateStream(&stream);
  if (FAILED(hr)) {
    cleanup();
    *error = "WIC stream creation failed";
    return false;
  }
  hr = stream->InitializeFromFilename(path.wstring().c_str(), GENERIC_WRITE);
  if (FAILED(hr)) {
    cleanup();
    *error = "PNG output file could not be opened";
    return false;
  }
  hr = factory->CreateEncoder(GUID_ContainerFormatPng, nullptr, &encoder);
  if (FAILED(hr) || FAILED(encoder->Initialize(stream, WICBitmapEncoderNoCache))) {
    cleanup();
    *error = "PNG encoder initialization failed";
    return false;
  }
  hr = encoder->CreateNewFrame(&frame, nullptr);
  if (FAILED(hr) || FAILED(frame->Initialize(nullptr)) || FAILED(frame->SetSize(width, height))) {
    cleanup();
    *error = "PNG frame initialization failed";
    return false;
  }
  WICPixelFormatGUID format = GUID_WICPixelFormat32bppBGRA;
  hr = frame->SetPixelFormat(&format);
  if (FAILED(hr)) {
    cleanup();
    *error = "PNG encoder pixel format setup failed";
    return false;
  }
  std::vector<BYTE> encoded;
  UINT stride = 0;
  if (format == GUID_WICPixelFormat32bppBGRA) {
    stride = width * 4;
    encoded.resize(rgba.size());
    for (size_t i = 0; i + 3 < rgba.size(); i += 4) {
      encoded[i] = rgba[i + 2];
      encoded[i + 1] = rgba[i + 1];
      encoded[i + 2] = rgba[i];
      encoded[i + 3] = rgba[i + 3];
    }
  } else if (format == GUID_WICPixelFormat32bppPBGRA) {
    stride = width * 4;
    encoded.resize(rgba.size());
    for (size_t i = 0; i + 3 < rgba.size(); i += 4) {
      const unsigned int alpha = rgba[i + 3];
      encoded[i] = static_cast<BYTE>((static_cast<unsigned int>(rgba[i + 2]) * alpha + 127) / 255);
      encoded[i + 1] = static_cast<BYTE>((static_cast<unsigned int>(rgba[i + 1]) * alpha + 127) / 255);
      encoded[i + 2] = static_cast<BYTE>((static_cast<unsigned int>(rgba[i]) * alpha + 127) / 255);
      encoded[i + 3] = rgba[i + 3];
    }
  } else if (format == GUID_WICPixelFormat24bppBGR) {
    stride = width * 3;
    encoded.resize(static_cast<size_t>(stride) * height);
    for (size_t src = 0, dst = 0; src + 3 < rgba.size() && dst + 2 < encoded.size(); src += 4, dst += 3) {
      encoded[dst] = rgba[src + 2];
      encoded[dst + 1] = rgba[src + 1];
      encoded[dst + 2] = rgba[src];
    }
  } else {
    cleanup();
    *error = "PNG encoder does not support required pixel format";
    return false;
  }
  const UINT size = stride * height;
  hr = frame->WritePixels(height, stride, size, encoded.data());
  if (FAILED(hr) || FAILED(frame->Commit()) || FAILED(encoder->Commit())) {
    cleanup();
    *error = "PNG writing failed";
    return false;
  }
  cleanup();
  return true;
}

std::vector<SenderInfo> sorted_senders(SPOUTHANDLE spout) {
  auto sender_names = spout->GetSenderList();
  std::sort(sender_names.begin(), sender_names.end());
  std::vector<SenderInfo> senders;
  senders.reserve(sender_names.size());
  for (const auto &name : sender_names) {
    SenderInfo info;
    info.name = name;
    unsigned int width = 0;
    unsigned int height = 0;
    HANDLE handle = nullptr;
    DWORD format = 0;
    char host_path[MAX_PATH] = {};
    spout->GetSenderInfo(name.c_str(), width, height, handle, format);
    spout->GetHostPath(name.c_str(), host_path, MAX_PATH);
    info.width = width;
    info.height = height;
    info.host_path = host_path;
    senders.push_back(std::move(info));
  }
  return senders;
}

struct SenderSelection {
  std::string name;
  std::string code;
  std::string message;
  std::vector<SenderInfo> senders;
};

bool sender_matches_requested(const std::string &candidate, const std::string &requested) {
  return candidate == requested || lower_copy(candidate) == lower_copy(requested);
}

SenderSelection choose_sender(SPOUTHANDLE spout, const std::string &requested) {
  SenderSelection selection;
  selection.senders = sorted_senders(spout);
  if (!requested.empty()) {
    for (const auto &sender : selection.senders) {
      if (sender_matches_requested(sender.name, requested)) {
        selection.name = sender.name;
        return selection;
      }
    }
    selection.code = "sender_not_found";
    selection.message = "指定されたSpout senderが見つかりません。候補一覧を確認してください。";
    return selection;
  }
  if (selection.senders.empty()) {
    selection.code = "sender_not_found";
    selection.message = "Spout senderがありません。VRChatでStream Cameraを起動してください。";
    return selection;
  }
  if (selection.senders.size() == 1) {
    selection.name = selection.senders[0].name;
    return selection;
  }
  std::vector<SenderInfo> auto_candidates;
  for (const auto &sender : selection.senders) {
    std::string lower = lower_copy(sender.name);
    if (lower.find("vrchat") != std::string::npos || lower.find("stream") != std::string::npos) {
      auto_candidates.push_back(sender);
    }
  }
  if (auto_candidates.size() == 1) {
    selection.name = auto_candidates[0].name;
    return selection;
  }
  selection.code = "sender_ambiguous";
  selection.message = "複数のSpout senderがあり自動選択できません。sender名を選択してください。";
  if (!auto_candidates.empty()) {
    selection.senders = std::move(auto_candidates);
  }
  return selection;
}

SenderSelection wait_for_sender(SPOUTHANDLE spout, const std::string &requested,
                                std::chrono::steady_clock::time_point deadline) {
  SenderSelection selection;
  while (std::chrono::steady_clock::now() < deadline) {
    selection = choose_sender(spout, requested);
    if (!selection.name.empty()) {
      return selection;
    }
    if (selection.code != "sender_not_found") {
      return selection;
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
  }
  selection = choose_sender(spout, requested);
  if (!selection.name.empty() || selection.code != "sender_not_found") {
    return selection;
  }
  if (requested.empty()) {
    selection.message = "Spout senderがありません。VRChatでStream Cameraを起動し、SpoutがONになるまで待ちましたが検出できませんでした。";
  } else {
    selection.message = "指定されたSpout senderが見つかりません。VRChatでStream Cameraを起動し、sender名を確認してください。";
  }
  return selection;
}

int list_senders(SPOUTHANDLE spout) {
  auto senders = sorted_senders(spout);
  std::cout << "{\"ok\":true,\"senders\":[";
  for (size_t i = 0; i < senders.size(); ++i) {
    if (i > 0) {
      std::cout << ",";
    }
    std::cout << "{\"name\":\"" << json_escape(senders[i].name) << "\",\"width\":"
              << senders[i].width << ",\"height\":" << senders[i].height
              << ",\"hostPath\":\"" << json_escape(senders[i].host_path) << "\"}";
  }
  std::cout << "]}\n";
  return 0;
}

int capture(SPOUTHANDLE spout, const Options &options) {
  auto deadline = std::chrono::steady_clock::now() + std::chrono::milliseconds(options.timeout_ms);
  SenderSelection selection = wait_for_sender(spout, options.sender, deadline);
  if (selection.name.empty()) {
    print_error(selection.code.empty() ? "sender_not_selected" : selection.code,
                selection.message, selection.senders);
    return 2;
  }
  spout->SetReceiverName(selection.name.c_str());
  std::vector<unsigned char> pixels;
  unsigned int width = 0;
  unsigned int height = 0;
  HANDLE handle = nullptr;
  DWORD format = 0;
  bool received_any = false;
  bool received_valid = false;
  FrameStats last_stats;
  while (std::chrono::steady_clock::now() < deadline) {
    spout->GetSenderInfo(selection.name.c_str(), width, height, handle, format);
    if (width == 0 || height == 0) {
      std::this_thread::sleep_for(std::chrono::milliseconds(30));
      continue;
    }
    pixels.assign(static_cast<size_t>(width) * static_cast<size_t>(height) * 4, 0);
    if (spout->ReceiveImage(pixels.data(), GL_RGBA, false, 0)) {
      received_any = true;
      last_stats = analyze_rgba_frame(pixels, width, height);
      if (!is_blank_frame(last_stats)) {
        force_opaque_alpha(pixels);
        received_valid = true;
        break;
      }
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(30));
  }
  if (!received_any) {
    print_error("capture_timeout",
                "Spoutフレームを取得できませんでした。VRChat Stream Cameraとsender設定を確認してください。",
                selection.senders);
    return 3;
  }
  if (!received_valid) {
    print_error("capture_blank_frame",
                "Spoutフレームは取得できましたが、timeout内に有効な映像になりませんでした。VRChat Stream Cameraの映像が表示されているか確認してください。",
                selection.senders);
    return 3;
  }
  std::error_code ec;
  if (!options.output.parent_path().empty()) {
    std::filesystem::create_directories(options.output.parent_path(), ec);
    if (ec) {
      print_error("output_directory_error", "出力フォルダを作成できません: " + ec.message());
      return 4;
    }
  }
  std::string write_error;
  if (!write_png_wic(options.output, width, height, pixels, &write_error)) {
    print_error("png_write_error", write_error);
    return 5;
  }
  auto now = std::chrono::system_clock::now();
  auto seconds = std::chrono::time_point_cast<std::chrono::seconds>(now);
  std::time_t t = std::chrono::system_clock::to_time_t(seconds);
  std::tm utc = {};
  gmtime_s(&utc, &t);
  char timestamp[32] = {};
  std::strftime(timestamp, sizeof(timestamp), "%Y-%m-%dT%H:%M:%SZ", &utc);
  std::cout << "{\"ok\":true,\"senderName\":\"" << json_escape(selection.name)
            << "\",\"width\":" << width << ",\"height\":" << height
            << ",\"frame\":" << spout->GetSenderFrame()
            << ",\"capturedAt\":\"" << timestamp << "\",\"outputPath\":\""
            << json_escape(options.output.u8string()) << "\"}\n";
  return 0;
}

} // namespace

int main(int argc, char **argv) {
  Options options;
  std::string parse_error;
  if (!parse_args(argc, argv, &options, &parse_error)) {
    print_error("invalid_arguments", parse_error);
    return 64;
  }
  SPOUTHANDLE spout = GetSpout();
  if (!spout) {
    print_error("spout_init_failed", "Spoutを初期化できませんでした。");
    return 1;
  }
  int rc = 0;
  if (options.list_senders) {
    rc = list_senders(spout);
  } else {
    rc = capture(spout, options);
  }
  spout->Release();
  return rc;
}
