#include <windows.h>
#include <wincodec.h>

#include <chrono>
#include <algorithm>
#include <cctype>
#include <cstddef>
#include <cstdint>
#include <ctime>
#include <filesystem>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <cmath>
#include <memory>
#include <sstream>
#include <string>
#include <utility>
#include <thread>
#include <vector>

#include "SpoutLibrary.h"
#include "capture_logic.h"

using spout_capture::CaptureFrameState;
using spout_capture::CaptureObservation;
using spout_capture::FrameStats;
using spout_capture::analyze_rgba_frame;
using spout_capture::capture_frame_state_code;
using spout_capture::capture_frame_state_message;
using spout_capture::classify_capture_frame_state;
using spout_capture::format_frame_stats;
using spout_capture::is_blank_frame;

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
  bool diagnose = false;
  bool show_version = false;
  std::string sender;
  std::filesystem::path output;
  std::filesystem::path debug_dir;
  int timeout_ms = 10000;
  int duration_ms = 10000;
  int interval_ms = 250;
  int debug_frames = 0;
  bool debug_frames_set = false;
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
            << "spout-capture --capture [--sender name] --output file.png --timeout-ms 10000\n"
            << "              [--debug-dir directory] [--debug-frames 8]\n"
            << "spout-capture --diagnose --debug-dir directory [--sender name]\n"
            << "              [--duration-ms 10000] [--interval-ms 250] [--debug-frames 8]\n";
}

void print_version() {
  std::cout << "{\"ok\":true,\"name\":\"spout-capture\",\"version\":\""
            << json_escape(helper_version()) << "\"}\n";
}

void write_senders_json(std::ostream &out, const std::vector<SenderInfo> &senders) {
  out << "[";
  for (size_t i = 0; i < senders.size(); ++i) {
    if (i > 0) {
      out << ",";
    }
    out << "{\"name\":\"" << json_escape(senders[i].name) << "\",\"width\":"
        << senders[i].width << ",\"height\":" << senders[i].height
        << ",\"hostPath\":\"" << json_escape(senders[i].host_path) << "\"}";
  }
  out << "]";
}

void print_error(const std::string &code, const std::string &message, const std::vector<SenderInfo> &senders = {}) {
  std::cout << "{\"ok\":false,\"code\":\"" << json_escape(code)
            << "\",\"message\":\"" << json_escape(message) << "\"";
  if (!senders.empty()) {
    std::cout << ",\"senders\":";
    write_senders_json(std::cout, senders);
  }
  std::cout << "}\n";
}

void print_capture_error(const std::string &code, const std::string &message, const std::string &sender_name,
                         unsigned int width, unsigned int height, const CaptureObservation &observation,
                         const FrameStats &stats, const std::vector<SenderInfo> &senders = {},
                         const std::string &debug_dir = {}, int debug_frames = 0) {
  std::cout << "{\"ok\":false,\"code\":\"" << json_escape(code)
            << "\",\"message\":\"" << json_escape(message) << "\""
            << ",\"senderName\":\"" << json_escape(sender_name) << "\""
            << ",\"width\":" << width << ",\"height\":" << height
            << ",\"frame\":" << observation.last_sender_frame
            << ",\"frameState\":\"" << capture_frame_state_code(classify_capture_frame_state(observation)) << "\""
            << ",\"receiveAttempts\":" << observation.receive_attempts
            << ",\"receiveSuccesses\":" << observation.receive_successes
            << ",\"firstFrame\":" << observation.first_sender_frame
            << ",\"lastReceivedFrame\":" << observation.last_received_frame
            << ",\"frameStats\":{\"samples\":" << stats.samples
            << ",\"mean\":" << stats.mean
            << ",\"stddev\":" << stats.stddev
            << ",\"nearWhiteRatio\":" << stats.near_white_ratio
            << ",\"nearBlackRatio\":" << stats.near_black_ratio
            << ",\"transparentRatio\":" << stats.transparent_ratio << "}"
            << ",\"frameStatsText\":\"" << json_escape(format_frame_stats(stats)) << "\"";
  if (!debug_dir.empty()) {
    std::cout << ",\"debugDir\":\"" << json_escape(debug_dir) << "\",\"debugFrames\":" << debug_frames;
  }
  if (!senders.empty()) {
    std::cout << ",\"senders\":";
    write_senders_json(std::cout, senders);
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
    } else if (arg == "--diagnose") {
      options->diagnose = true;
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
    } else if (arg == "--debug-dir") {
      if (++i >= argc) {
        *error = "--debug-dir requires a value";
        return false;
      }
      options->debug_dir = std::filesystem::path(argv[i]);
    } else if (arg == "--debug-frames") {
      if (++i >= argc || !parse_int(argv[i], &options->debug_frames)) {
        *error = "--debug-frames requires an integer value";
        return false;
      }
      options->debug_frames_set = true;
    } else if (arg == "--timeout-ms") {
      if (++i >= argc || !parse_int(argv[i], &options->timeout_ms)) {
        *error = "--timeout-ms requires an integer value";
        return false;
      }
    } else if (arg == "--duration-ms") {
      if (++i >= argc || !parse_int(argv[i], &options->duration_ms)) {
        *error = "--duration-ms requires an integer value";
        return false;
      }
    } else if (arg == "--interval-ms") {
      if (++i >= argc || !parse_int(argv[i], &options->interval_ms)) {
        *error = "--interval-ms requires an integer value";
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
  int modes = 0;
  if (options->list_senders) {
    modes++;
  }
  if (options->capture) {
    modes++;
  }
  if (options->diagnose) {
    modes++;
  }
  if (modes != 1) {
    *error = "specify exactly one of --list-senders, --capture, or --diagnose";
    return false;
  }
  if (options->capture && options->output.empty()) {
    *error = "--capture requires --output";
    return false;
  }
  if (options->diagnose && options->debug_dir.empty()) {
    *error = "--diagnose requires --debug-dir";
    return false;
  }
  if (options->timeout_ms < 100) {
    options->timeout_ms = 100;
  }
  if (options->duration_ms < 100) {
    options->duration_ms = 100;
  }
  if (options->duration_ms > 120000) {
    options->duration_ms = 120000;
  }
  if (options->interval_ms < 30) {
    options->interval_ms = 30;
  }
  if (options->interval_ms > 5000) {
    options->interval_ms = 5000;
  }
  if (!options->debug_dir.empty()) {
    if (options->debug_frames <= 0) {
      options->debug_frames = 8;
    }
    if (options->debug_frames > 120) {
      options->debug_frames = 120;
    }
  } else if (options->debug_frames_set) {
    *error = "--debug-frames requires --debug-dir";
    return false;
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

void force_opaque_alpha(std::vector<unsigned char> &rgba) {
  for (size_t i = 3; i < rgba.size(); i += 4) {
    rgba[i] = 255;
  }
}

std::string timestamp_utc() {
  auto now = std::chrono::system_clock::now();
  auto seconds = std::chrono::time_point_cast<std::chrono::seconds>(now);
  std::time_t t = std::chrono::system_clock::to_time_t(seconds);
  std::tm utc = {};
  gmtime_s(&utc, &t);
  char timestamp[32] = {};
  std::strftime(timestamp, sizeof(timestamp), "%Y-%m-%dT%H:%M:%SZ", &utc);
  return timestamp;
}

std::string timestamp_file_utc() {
  auto now = std::chrono::system_clock::now();
  auto seconds = std::chrono::time_point_cast<std::chrono::seconds>(now);
  std::time_t t = std::chrono::system_clock::to_time_t(seconds);
  std::tm utc = {};
  gmtime_s(&utc, &t);
  char timestamp[32] = {};
  std::strftime(timestamp, sizeof(timestamp), "%Y%m%dT%H%M%SZ", &utc);
  return timestamp;
}

void write_frame_stats_json(std::ostream &out, const FrameStats &stats) {
  out << "{\"samples\":" << stats.samples
      << ",\"mean\":" << stats.mean
      << ",\"stddev\":" << stats.stddev
      << ",\"nearWhiteRatio\":" << stats.near_white_ratio
      << ",\"nearBlackRatio\":" << stats.near_black_ratio
      << ",\"transparentRatio\":" << stats.transparent_ratio << "}";
}

struct DebugRecorder {
  bool enabled = false;
  std::filesystem::path dir;
  std::ofstream log;
  std::ofstream frames_jsonl;
  std::string session_id;
  int max_frames = 0;
  int dumped_frames = 0;

  void start(const Options &options, std::string *error) {
    if (options.debug_dir.empty()) {
      return;
    }
    dir = options.debug_dir;
    session_id = timestamp_file_utc() + "_pid" + std::to_string(GetCurrentProcessId());
    max_frames = options.debug_frames;
    std::error_code ec;
    std::filesystem::create_directories(dir, ec);
    if (ec) {
      *error = "debug directory could not be created: " + ec.message();
      return;
    }
    log.open(dir / "spout-capture-debug.log", std::ios::app);
    frames_jsonl.open(dir / "frames.jsonl", std::ios::app);
    if (!log || !frames_jsonl) {
      *error = "debug log files could not be opened";
      return;
    }
    enabled = true;
    log_event("debug recording started dir=\"" + dir.u8string() + "\" session=\"" + session_id +
              "\" max_frames=" + std::to_string(max_frames));
  }

  void log_event(const std::string &message) {
    if (!enabled || !log) {
      return;
    }
    log << timestamp_utc() << " " << message << "\n";
    log.flush();
  }

  void dump_frame(const std::vector<unsigned char> &rgba, unsigned int width, unsigned int height,
                  int64_t sender_frame, const FrameStats &stats) {
    if (!enabled || dumped_frames >= max_frames || width == 0 || height == 0 || rgba.empty()) {
      return;
    }
    const int next_index = dumped_frames + 1;
    std::ostringstream base;
    base << session_id << "_frame_" << std::setw(6) << std::setfill('0') << next_index;
    const std::filesystem::path raw_path = dir / (base.str() + ".rgba");
    const std::filesystem::path json_path = dir / (base.str() + ".json");
    std::ofstream raw(raw_path, std::ios::binary);
    if (!raw) {
      log_event("debug frame raw file open failed path=\"" + raw_path.u8string() + "\"");
      return;
    }
    raw.write(reinterpret_cast<const char *>(rgba.data()), static_cast<std::streamsize>(rgba.size()));
    if (!raw) {
      log_event("debug frame raw file write failed path=\"" + raw_path.u8string() + "\"");
      return;
    }
    dumped_frames = next_index;
    std::ostringstream metadata;
    metadata << "{\"capturedAt\":\"" << timestamp_utc()
             << "\",\"session\":\"" << json_escape(session_id) << "\""
             << ",\"index\":" << dumped_frames
             << ",\"senderFrame\":" << sender_frame
             << ",\"width\":" << width
             << ",\"height\":" << height
             << ",\"format\":\"rgba8\""
             << ",\"rawPath\":\"" << json_escape(raw_path.u8string()) << "\""
             << ",\"frameStats\":";
    write_frame_stats_json(metadata, stats);
    metadata << "}";
    std::ofstream json(json_path);
    if (json) {
      json << metadata.str() << "\n";
    }
    if (frames_jsonl) {
      frames_jsonl << metadata.str() << "\n";
      frames_jsonl.flush();
    }
    log_event("debug frame dumped index=" + std::to_string(dumped_frames) +
              " sender_frame=" + std::to_string(sender_frame) +
              " stats=\"" + format_frame_stats(stats) + "\"");
  }
};

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
  std::cout << "{\"ok\":true,\"senders\":";
  write_senders_json(std::cout, senders);
  std::cout << "}\n";
  return 0;
}

struct DiagnoseSummary {
  int samples = 0;
  int sender_samples = 0;
  int empty_sender_samples = 0;
  int ambiguous_sender_samples = 0;
  int receive_attempts = 0;
  int receive_successes = 0;
  int valid_frames = 0;
  int blank_frames = 0;
  int64_t first_frame = -1;
  int64_t last_frame = -1;
  bool frame_advanced = false;
  std::string sender_name;
  unsigned int width = 0;
  unsigned int height = 0;
  FrameStats last_stats;
};

std::string diagnose_code(const DiagnoseSummary &summary) {
  if (summary.valid_frames > 0) {
    return "diagnose_success";
  }
  if (summary.receive_successes > 0) {
    return "diagnose_blank_frame";
  }
  if (summary.frame_advanced) {
    return "diagnose_receive_stalled";
  }
  if (summary.sender_samples > 0) {
    return "diagnose_no_new_frame";
  }
  if (summary.ambiguous_sender_samples > 0) {
    return "diagnose_sender_ambiguous";
  }
  return "diagnose_no_sender";
}

std::string diagnose_message(const std::string &code) {
  if (code == "diagnose_success") {
    return "診断中に有効なSpout映像フレームを確認しました。";
  }
  if (code == "diagnose_blank_frame") {
    return "Spoutフレームは受信できましたが、診断中は有効な映像になりませんでした。";
  }
  if (code == "diagnose_receive_stalled") {
    return "Spout senderのフレーム番号は進みましたが、画像を受信できませんでした。";
  }
  if (code == "diagnose_no_new_frame") {
    return "Spout senderは見つかりましたが、診断中に新しいフレームを確認できませんでした。";
  }
  if (code == "diagnose_sender_ambiguous") {
    return "複数のSpout senderがあり、自動選択できませんでした。";
  }
  return "診断中にSpout senderを検出できませんでした。";
}

void write_diagnose_summary_json(std::ostream &out, const DiagnoseSummary &summary,
                                 const Options &options, const std::string &code,
                                 const std::filesystem::path &timeline_path,
                                 const std::filesystem::path &summary_path,
                                 int debug_frames) {
  out << "{\"ok\":true"
      << ",\"code\":\"" << json_escape(code) << "\""
      << ",\"message\":\"" << json_escape(diagnose_message(code)) << "\""
      << ",\"debugDir\":\"" << json_escape(options.debug_dir.u8string()) << "\""
      << ",\"durationMs\":" << options.duration_ms
      << ",\"intervalMs\":" << options.interval_ms
      << ",\"samples\":" << summary.samples
      << ",\"senderName\":\"" << json_escape(summary.sender_name) << "\""
      << ",\"senderSamples\":" << summary.sender_samples
      << ",\"emptySenderSamples\":" << summary.empty_sender_samples
      << ",\"ambiguousSenderSamples\":" << summary.ambiguous_sender_samples
      << ",\"width\":" << summary.width
      << ",\"height\":" << summary.height
      << ",\"firstFrame\":" << summary.first_frame
      << ",\"lastFrame\":" << summary.last_frame
      << ",\"frameAdvanced\":" << (summary.frame_advanced ? "true" : "false")
      << ",\"receiveAttempts\":" << summary.receive_attempts
      << ",\"receiveSuccesses\":" << summary.receive_successes
      << ",\"validFrames\":" << summary.valid_frames
      << ",\"blankFrames\":" << summary.blank_frames
      << ",\"lastFrameStats\":";
  write_frame_stats_json(out, summary.last_stats);
  out << ",\"timelinePath\":\"" << json_escape(timeline_path.u8string()) << "\""
      << ",\"summaryPath\":\"" << json_escape(summary_path.u8string()) << "\""
      << ",\"debugFrames\":" << debug_frames
      << "}";
}

int diagnose(SPOUTHANDLE spout, const Options &options) {
  DebugRecorder debug;
  std::string debug_error;
  debug.start(options, &debug_error);
  if (!debug_error.empty()) {
    print_error("debug_record_error", "Spout debug録画の準備に失敗しました: " + debug_error);
    return 6;
  }
  const std::filesystem::path timeline_path = options.debug_dir / "diagnose.jsonl";
  const std::filesystem::path summary_path = options.debug_dir / "diagnose-summary.json";
  std::ofstream timeline(timeline_path, std::ios::app);
  if (!timeline) {
    print_error("debug_record_error", "Spout diagnose timelineを開けませんでした: " + timeline_path.u8string());
    return 6;
  }

  DiagnoseSummary summary;
  std::vector<unsigned char> pixels;
  unsigned int width = 0;
  unsigned int height = 0;
  HANDLE handle = nullptr;
  DWORD format = 0;
  std::string receiver_name;
  const auto started = std::chrono::steady_clock::now();
  const auto deadline = started + std::chrono::milliseconds(options.duration_ms);

  while (std::chrono::steady_clock::now() < deadline) {
    const auto now = std::chrono::steady_clock::now();
    const int64_t elapsed_ms = std::chrono::duration_cast<std::chrono::milliseconds>(now - started).count();
    SenderSelection selection = choose_sender(spout, options.sender);
    summary.samples++;
    if (selection.senders.empty()) {
      summary.empty_sender_samples++;
    }
    if (selection.code == "sender_ambiguous") {
      summary.ambiguous_sender_samples++;
    }

    bool receive_attempted = false;
    bool receive_ok = false;
    FrameStats stats;
    int64_t sender_frame = -1;
    std::string frame_state = "no_sender";
    width = 0;
    height = 0;
    handle = nullptr;
    format = 0;
    if (!selection.name.empty()) {
      summary.sender_samples++;
      if (receiver_name != selection.name) {
        spout->SetReceiverName(selection.name.c_str());
        receiver_name = selection.name;
        debug.log_event("diagnose sender selected name=\"" + selection.name + "\"");
      }
      summary.sender_name = selection.name;
      sender_frame = spout->GetSenderFrame();
      if (summary.first_frame < 0) {
        summary.first_frame = sender_frame;
      } else if (sender_frame != summary.last_frame) {
        summary.frame_advanced = true;
      }
      summary.last_frame = sender_frame;

      spout->GetSenderInfo(selection.name.c_str(), width, height, handle, format);
      summary.width = width;
      summary.height = height;
      if (width > 0 && height > 0) {
        pixels.assign(static_cast<size_t>(width) * static_cast<size_t>(height) * 4, 0);
        receive_attempted = true;
        summary.receive_attempts++;
        if (spout->ReceiveImage(pixels.data(), GL_RGBA, false, 0)) {
          receive_ok = true;
          summary.receive_successes++;
          stats = analyze_rgba_frame(pixels, width, height);
          summary.last_stats = stats;
          debug.dump_frame(pixels, width, height, sender_frame, stats);
          if (is_blank_frame(stats)) {
            summary.blank_frames++;
            frame_state = "blank";
          } else {
            summary.valid_frames++;
            frame_state = "valid";
          }
        } else {
          frame_state = "receive_failed";
        }
      } else {
        frame_state = "zero_size_sender";
      }
    } else if (!selection.code.empty()) {
      frame_state = selection.code;
    }

    timeline << "{\"capturedAt\":\"" << timestamp_utc()
             << "\",\"elapsedMs\":" << elapsed_ms
             << ",\"senders\":";
    write_senders_json(timeline, selection.senders);
    timeline << ",\"selectedSender\":\"" << json_escape(selection.name) << "\""
             << ",\"senderFound\":" << (!selection.name.empty() ? "true" : "false")
             << ",\"width\":" << width
             << ",\"height\":" << height
             << ",\"senderFrame\":" << sender_frame
             << ",\"frameAdvanced\":" << (summary.frame_advanced ? "true" : "false")
             << ",\"receiveAttempted\":" << (receive_attempted ? "true" : "false")
             << ",\"receiveOK\":" << (receive_ok ? "true" : "false")
             << ",\"frameState\":\"" << json_escape(frame_state) << "\""
             << ",\"frameStats\":";
    write_frame_stats_json(timeline, stats);
    timeline << "}\n";
    timeline.flush();

    const auto next_tick = now + std::chrono::milliseconds(options.interval_ms);
    if (next_tick < deadline) {
      std::this_thread::sleep_until(next_tick);
    } else {
      break;
    }
  }

  const std::string code = diagnose_code(summary);
  std::ofstream summary_file(summary_path);
  if (summary_file) {
    write_diagnose_summary_json(summary_file, summary, options, code, timeline_path, summary_path, debug.dumped_frames);
    summary_file << "\n";
  } else {
    debug.log_event("diagnose summary file open failed path=\"" + summary_path.u8string() + "\"");
  }
  debug.log_event("diagnose complete code=\"" + code +
                  "\" samples=" + std::to_string(summary.samples) +
                  " sender_samples=" + std::to_string(summary.sender_samples) +
                  " receive_successes=" + std::to_string(summary.receive_successes) +
                  " valid_frames=" + std::to_string(summary.valid_frames) +
                  " blank_frames=" + std::to_string(summary.blank_frames));
  write_diagnose_summary_json(std::cout, summary, options, code, timeline_path, summary_path, debug.dumped_frames);
  std::cout << "\n";
  return 0;
}

int capture(SPOUTHANDLE spout, const Options &options) {
  DebugRecorder debug;
  std::string debug_error;
  debug.start(options, &debug_error);
  if (!debug_error.empty()) {
    print_error("debug_record_error", "Spout debug録画の準備に失敗しました: " + debug_error);
    return 6;
  }
  auto deadline = std::chrono::steady_clock::now() + std::chrono::milliseconds(options.timeout_ms);
  SenderSelection selection = wait_for_sender(spout, options.sender, deadline);
  if (selection.name.empty()) {
    debug.log_event("sender selection failed code=\"" + selection.code + "\" message=\"" + selection.message + "\"");
    print_error(selection.code.empty() ? "sender_not_selected" : selection.code,
                selection.message, selection.senders);
    return 2;
  }
  debug.log_event("sender selected name=\"" + selection.name + "\"");
  spout->SetReceiverName(selection.name.c_str());
  std::vector<unsigned char> pixels;
  unsigned int width = 0;
  unsigned int height = 0;
  HANDLE handle = nullptr;
  DWORD format = 0;
  CaptureObservation observation;
  while (std::chrono::steady_clock::now() < deadline) {
    const int64_t current_frame = spout->GetSenderFrame();
    if (observation.first_sender_frame < 0) {
      observation.first_sender_frame = current_frame;
    } else if (current_frame != observation.last_sender_frame) {
      observation.saw_frame_advance = true;
    }
    observation.last_sender_frame = current_frame;

    spout->GetSenderInfo(selection.name.c_str(), width, height, handle, format);
    if (width == 0 || height == 0) {
      debug.log_event("sender info has zero size sender_frame=" + std::to_string(current_frame));
      std::this_thread::sleep_for(std::chrono::milliseconds(30));
      continue;
    }

    pixels.assign(static_cast<size_t>(width) * static_cast<size_t>(height) * 4, 0);
    observation.receive_attempts++;
    if (spout->ReceiveImage(pixels.data(), GL_RGBA, false, 0)) {
      observation.saw_receive_success = true;
      observation.receive_successes++;
      observation.last_received_frame = current_frame;
      observation.last_stats = analyze_rgba_frame(pixels, width, height);
      debug.dump_frame(pixels, width, height, current_frame, observation.last_stats);
      if (!is_blank_frame(observation.last_stats)) {
        observation.saw_non_blank_frame = true;
        debug.log_event("valid non blank frame found sender_frame=" + std::to_string(current_frame) +
                        " stats=\"" + format_frame_stats(observation.last_stats) + "\"");
        force_opaque_alpha(pixels);
        break;
      }
    } else {
      debug.log_event("ReceiveImage returned false sender_frame=" + std::to_string(current_frame));
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(30));
  }

  const CaptureFrameState state = classify_capture_frame_state(observation);
  if (state != CaptureFrameState::ValidFrame) {
    debug.log_event("capture failed state=\"" + std::string(capture_frame_state_code(state)) +
                    "\" attempts=" + std::to_string(observation.receive_attempts) +
                    " successes=" + std::to_string(observation.receive_successes) +
                    " first_frame=" + std::to_string(observation.first_sender_frame) +
                    " last_frame=" + std::to_string(observation.last_sender_frame) +
                    " dumped_frames=" + std::to_string(debug.dumped_frames));
    print_capture_error(capture_frame_state_code(state), capture_frame_state_message(state),
                        selection.name, width, height, observation, observation.last_stats,
                        selection.senders, debug.enabled ? debug.dir.u8string() : "", debug.dumped_frames);
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
    debug.log_event("png write failed error=\"" + write_error + "\"");
    print_error("png_write_error", write_error);
    return 5;
  }
  const std::string timestamp = timestamp_utc();
  debug.log_event("capture success output=\"" + options.output.u8string() + "\" dumped_frames=" + std::to_string(debug.dumped_frames));
  std::cout << "{\"ok\":true,\"senderName\":\"" << json_escape(selection.name)
            << "\",\"width\":" << width << ",\"height\":" << height
            << ",\"frame\":" << observation.last_sender_frame
            << ",\"capturedAt\":\"" << timestamp << "\",\"outputPath\":\""
            << json_escape(options.output.u8string()) << "\"";
  if (debug.enabled) {
    std::cout << ",\"debugDir\":\"" << json_escape(debug.dir.u8string())
              << "\",\"debugFrames\":" << debug.dumped_frames;
  }
  std::cout << "}\n";
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
  } else if (options.diagnose) {
    rc = diagnose(spout, options);
  } else {
    rc = capture(spout, options);
  }
  spout->Release();
  return rc;
}
