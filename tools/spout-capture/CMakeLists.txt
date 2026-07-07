cmake_minimum_required(VERSION 3.21)
project(spout_capture LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

enable_testing()

add_executable(spout-capture-logic-test capture_logic_test.cpp)
target_include_directories(spout-capture-logic-test PRIVATE ${CMAKE_CURRENT_SOURCE_DIR})
add_test(NAME spout-capture-logic-test COMMAND spout-capture-logic-test)

if (WIN32)
  include(FetchContent)

  set(SPOUT2_REVISION "fc63ac2c16918ffa1d3384cddfba2bbc70650bb4")
  set(SPOUT2_UPSTREAM_TAG "2.007.015")
  set(SPOUT2_ARCHIVE_URL "https://github.com/leadedge/Spout2/archive/${SPOUT2_REVISION}.tar.gz")
  set(SPOUT2_ARCHIVE_SHA256 "087730642decf46d71fa1a96c1f63f4f6a99feb5e82e9b911b6b70c65d8a6305")

  FetchContent_Declare(
    Spout2
    URL ${SPOUT2_ARCHIVE_URL}
    URL_HASH SHA256=${SPOUT2_ARCHIVE_SHA256}
    DOWNLOAD_EXTRACT_TIMESTAMP FALSE
  )
  FetchContent_MakeAvailable(Spout2)

  add_executable(spout-capture main.cpp)
  target_link_libraries(spout-capture PRIVATE SpoutLibrary windowscodecs ole32)
  target_include_directories(spout-capture PRIVATE
    ${spout2_SOURCE_DIR}/SPOUTSDK/SpoutLibrary
  )
  add_custom_command(TARGET spout-capture POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy_if_different
      $<TARGET_FILE:SpoutLibrary>
      $<TARGET_FILE_DIR:spout-capture>
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
endif()
