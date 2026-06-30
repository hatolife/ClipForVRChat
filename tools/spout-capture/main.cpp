#include <windows.h>
#include <wincodec.h>

#include <chrono>
#include <algorithm>
#include <cctype>
#include <cstdint>
#include <ctime>
#include <filesystem>
#include <iostream>
#include <memory>
#include <sstream>
#include <string>
#include <thread>
#include <vector>

#include "SpoutLibrary.h"

namespace {

struct Options {
  bool list_senders = false;
  bool capture = false;
  std::string sender;
  std::filesystem::path output;
  int timeout_ms = 10000;
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

void print_error(const std::string &code, const std::string &message) {
  std::cout << "{\"ok\":false,\"code\":\"" << json_escape(code)
            << "\",\"message\":\"" << json_escape(message) << "\"}\n";
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
      std::cout << "spout-capture --list-senders\n"
                << "spout-capture --capture [--sender name] --output file.png --timeout-ms 10000\n";
      std::exit(0);
    } else {
      *error = "unknown argument: " + arg;
      return false;
    }
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
  WICPixelFormatGUID format = GUID_WICPixelFormat32bppRGBA;
  hr = frame->SetPixelFormat(&format);
  if (FAILED(hr) || format != GUID_WICPixelFormat32bppRGBA) {
    cleanup();
    *error = "PNG encoder does not support RGBA";
    return false;
  }
  const UINT stride = width * 4;
  const UINT size = stride * height;
  hr = frame->WritePixels(height, stride, size, const_cast<BYTE *>(rgba.data()));
  if (FAILED(hr) || FAILED(frame->Commit()) || FAILED(encoder->Commit())) {
    cleanup();
    *error = "PNG writing failed";
    return false;
  }
  cleanup();
  return true;
}

std::vector<std::string> sorted_senders(SPOUTHANDLE spout) {
  auto senders = spout->GetSenderList();
  std::sort(senders.begin(), senders.end());
  return senders;
}

std::string choose_sender(SPOUTHANDLE spout, const std::string &requested, std::string *error) {
  if (!requested.empty()) {
    return requested;
  }
  auto senders = sorted_senders(spout);
  if (senders.empty()) {
    *error = "Spout senderがありません。VRChatでStream Cameraを起動してください。";
    return {};
  }
  if (senders.size() == 1) {
    return senders[0];
  }
  for (const auto &sender : senders) {
    std::string lower = sender;
    std::transform(lower.begin(), lower.end(), lower.begin(), [](unsigned char c) { return static_cast<char>(std::tolower(c)); });
    if (lower.find("vrchat") != std::string::npos || lower.find("stream") != std::string::npos) {
      return sender;
    }
  }
  *error = "複数のSpout senderがあり自動選択できません。sender名を選択してください。";
  return {};
}

int list_senders(SPOUTHANDLE spout) {
  auto senders = sorted_senders(spout);
  std::cout << "{\"ok\":true,\"senders\":[";
  for (size_t i = 0; i < senders.size(); ++i) {
    unsigned int width = 0;
    unsigned int height = 0;
    HANDLE handle = nullptr;
    DWORD format = 0;
    char host_path[MAX_PATH] = {};
    spout->GetSenderInfo(senders[i].c_str(), width, height, handle, format);
    spout->GetHostPath(senders[i].c_str(), host_path, MAX_PATH);
    if (i > 0) {
      std::cout << ",";
    }
    std::cout << "{\"name\":\"" << json_escape(senders[i]) << "\",\"width\":" << width
              << ",\"height\":" << height << ",\"hostPath\":\"" << json_escape(host_path) << "\"}";
  }
  std::cout << "]}\n";
  return 0;
}

int capture(SPOUTHANDLE spout, const Options &options) {
  std::string choose_error;
  std::string sender = choose_sender(spout, options.sender, &choose_error);
  if (sender.empty()) {
    print_error("sender_not_selected", choose_error);
    return 2;
  }
  spout->SetReceiverName(sender.c_str());
  auto deadline = std::chrono::steady_clock::now() + std::chrono::milliseconds(options.timeout_ms);
  std::vector<unsigned char> pixels;
  unsigned int width = 0;
  unsigned int height = 0;
  HANDLE handle = nullptr;
  DWORD format = 0;
  bool received = false;
  while (std::chrono::steady_clock::now() < deadline) {
    spout->GetSenderInfo(sender.c_str(), width, height, handle, format);
    if (width == 0 || height == 0) {
      std::this_thread::sleep_for(std::chrono::milliseconds(30));
      continue;
    }
    pixels.assign(static_cast<size_t>(width) * static_cast<size_t>(height) * 4, 0);
    if (spout->ReceiveImage(pixels.data(), GL_RGBA, false, 0)) {
      received = true;
      break;
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(30));
  }
  if (!received) {
    print_error("capture_timeout", "Spoutフレームを取得できませんでした。VRChat Stream Cameraとsender設定を確認してください。");
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
  std::cout << "{\"ok\":true,\"senderName\":\"" << json_escape(spout->GetSenderName())
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
