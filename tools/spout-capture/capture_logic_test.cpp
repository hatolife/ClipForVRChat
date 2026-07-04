#include "capture_logic.h"

#include <cstdlib>
#include <iostream>
#include <vector>

using spout_capture::CaptureFrameState;
using spout_capture::CaptureObservation;

namespace {

void expect(bool condition, const char *message) {
  if (!condition) {
    std::cerr << message << "\n";
    std::exit(1);
  }
}

void test_analyze_rgba_frame_tracks_transparency() {
  std::vector<unsigned char> rgba = {
      0, 0, 0, 0,
      0, 0, 0, 0,
      0, 0, 0, 0,
      0, 0, 0, 0,
  };
  auto stats = spout_capture::analyze_rgba_frame(rgba, 1, 1);
  expect(stats.samples == 1, "transparent frame should produce one sample");
  expect(stats.transparent_ratio == 1.0, "transparent frame should report transparency");
  expect(spout_capture::is_blank_frame(stats), "transparent frame should be blank");
}

void test_classify_capture_frame_state() {
  CaptureObservation no_new_frame;
  expect(spout_capture::classify_capture_frame_state(no_new_frame) == CaptureFrameState::NoNewFrame,
         "empty observation should classify as no new frame");

  CaptureObservation stalled;
  stalled.saw_frame_advance = true;
  expect(spout_capture::classify_capture_frame_state(stalled) == CaptureFrameState::ReceiveStalled,
         "frame advance without receive should classify as receive stalled");

  CaptureObservation blank;
  blank.saw_receive_success = true;
  expect(spout_capture::classify_capture_frame_state(blank) == CaptureFrameState::BlankFrame,
         "blank receive without progress should classify as blank frame");

  CaptureObservation valid;
  valid.saw_receive_success = true;
  valid.saw_frame_advance = true;
  valid.saw_non_blank_frame = true;
  expect(spout_capture::classify_capture_frame_state(valid) == CaptureFrameState::ValidFrame,
         "non blank receive should classify as valid frame");
}

} // namespace

int main() {
  test_analyze_rgba_frame_tracks_transparency();
  test_classify_capture_frame_state();
  return 0;
}
