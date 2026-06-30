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

if (MSVC)
  target_compile_options(spout-capture PRIVATE /W4 /permissive-)
endif()
