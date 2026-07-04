#pragma once

#include <algorithm>
#include <cmath>
#include <cstdint>
#include <cstddef>
#include <iomanip>
#include <sstream>
#include <string>
#include <vector>

namespace spout_capture {

struct FrameStats {
  int samples = 0;
  double mean = 0.0;
  double stddev = 0.0;
  double near_white_ratio = 0.0;
  double near_black_ratio = 0.0;
  double transparent_ratio = 0.0;
};

struct CaptureObservation {
  bool saw_receive_success = false;
  bool saw_frame_advance = false;
  bool saw_non_blank_frame = false;
  int receive_attempts = 0;
  int receive_successes = 0;
  int64_t first_sender_frame = -1;
  int64_t last_sender_frame = -1;
  int64_t last_received_frame = -1;
  FrameStats last_stats;
};

enum class CaptureFrameState {
  NoNewFrame,
  ReceiveStalled,
  BlankFrame,
  ValidFrame,
};

inline FrameStats analyze_rgba_frame(const std::vector<unsigned char> &rgba, unsigned int width, unsigned int height) {
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
  double transparent = 0.0;
  for (unsigned int y = 0; y < height; y += step_y) {
    for (unsigned int x = 0; x < width; x += step_x) {
      const size_t i = (static_cast<size_t>(y) * width + x) * 4;
      if (i + 3 >= rgba.size()) {
        continue;
      }
      const double r = static_cast<double>(rgba[i]);
      const double g = static_cast<double>(rgba[i + 1]);
      const double b = static_cast<double>(rgba[i + 2]);
      const double a = static_cast<double>(rgba[i + 3]);
      const double luma = 0.2126 * r + 0.7152 * g + 0.0722 * b;
      sum += luma;
      sum_sq += luma * luma;
      if (luma >= 250.0 && a >= 250.0) {
        near_white++;
      }
      if (luma <= 5.0 && a >= 250.0) {
        near_black++;
      }
      if (a <= 5.0) {
        transparent++;
      }
      stats.samples++;
    }
  }
  if (stats.samples == 0) {
    return stats;
  }
  const double samples = static_cast<double>(stats.samples);
  stats.mean = sum / samples;
  const double variance = (std::max)(0.0, (sum_sq / samples) - (stats.mean * stats.mean));
  stats.stddev = std::sqrt(variance);
  stats.near_white_ratio = near_white / samples;
  stats.near_black_ratio = near_black / samples;
  stats.transparent_ratio = transparent / samples;
  return stats;
}

inline bool is_blank_frame(const FrameStats &stats) {
  if (stats.samples == 0) {
    return true;
  }
  if (stats.transparent_ratio > 0.99) {
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

inline CaptureFrameState classify_capture_frame_state(const CaptureObservation &observation) {
  if (observation.saw_non_blank_frame) {
    return CaptureFrameState::ValidFrame;
  }
  if (observation.saw_receive_success) {
    return CaptureFrameState::BlankFrame;
  }
  if (observation.saw_frame_advance) {
    return CaptureFrameState::ReceiveStalled;
  }
  return CaptureFrameState::NoNewFrame;
}

inline const char *capture_frame_state_code(CaptureFrameState state) {
  switch (state) {
  case CaptureFrameState::NoNewFrame:
    return "capture_no_new_frame";
  case CaptureFrameState::ReceiveStalled:
    return "capture_receive_stalled";
  case CaptureFrameState::BlankFrame:
    return "capture_blank_frame";
  case CaptureFrameState::ValidFrame:
    return "capture_success";
  }
  return "capture_unknown";
}

inline const char *capture_frame_state_message(CaptureFrameState state) {
  switch (state) {
  case CaptureFrameState::NoNewFrame:
    return "Spout senderは見つかりましたが、timeout内に新しいフレームを受信できませんでした。VRChat Stream Cameraが更新されているか確認してください。";
  case CaptureFrameState::ReceiveStalled:
    return "Spout senderのフレーム番号は進みましたが、画像を受信できませんでした。VRChat Stream CameraとSpout受信状態を確認してください。";
  case CaptureFrameState::BlankFrame:
    return "Spoutフレームは取得できましたが、timeout内に有効な映像になりませんでした。VRChat Stream Cameraの映像が表示されているか確認してください。";
  case CaptureFrameState::ValidFrame:
    return "Spoutフレームを取得できました。";
  }
  return "Spout取得状態を判定できませんでした。";
}

inline std::string format_frame_stats(const FrameStats &stats) {
  std::ostringstream out;
  out << "samples=" << stats.samples
      << " mean=" << std::fixed << std::setprecision(2) << stats.mean
      << " stddev=" << std::fixed << std::setprecision(2) << stats.stddev
      << " near_white=" << std::fixed << std::setprecision(4) << stats.near_white_ratio
      << " near_black=" << std::fixed << std::setprecision(4) << stats.near_black_ratio
      << " transparent=" << std::fixed << std::setprecision(4) << stats.transparent_ratio;
  return out.str();
}

} // namespace spout_capture
