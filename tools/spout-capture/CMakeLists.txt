cmake_minimum_required(VERSION 3.21)
project(spout_capture LANGUAGES CXX)

include(FetchContent)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

FetchContent_Declare(
  Spout2
  GIT_REPOSITORY https://github.com/leadedge/Spout2.git
  GIT_TAG 2.007.015
)
FetchContent_MakeAvailable(Spout2)

add_executable(spout-capture main.cpp)
target_link_libraries(spout-capture PRIVATE SpoutLibrary_static windowscodecs ole32)
target_include_directories(spout-capture PRIVATE
  ${spout2_SOURCE_DIR}/SPOUTSDK/SpoutLibrary
)

execute_process(
  COMMAND git -C ${CMAKE_CURRENT_SOURCE_DIR} rev-parse --short=7 HEAD
  OUTPUT_VARIABLE SPOUT_CAPTURE_GIT_REVISION
  OUTPUT_STRIP_TRAILING_WHITESPACE
  ERROR_QUIET
)
if (NOT SPOUT_CAPTURE_GIT_REVISION)
  set(SPOUT_CAPTURE_GIT_REVISION "unknown")
endif()
target_compile_definitions(spout-capture PRIVATE
  SPOUT_CAPTURE_VERSION="v0.1.8-${SPOUT_CAPTURE_GIT_REVISION}"
)

if (MSVC)
  target_compile_options(spout-capture PRIVATE /W4 /permissive-)
endif()
