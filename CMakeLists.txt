cmake_minimum_required(VERSION 3.10)
set(CMAKE_CXX_STANDARD 17)

# project name
project(glitchfuck)

# pkg-config
find_package(PkgConfig)

include_directories(${CMAKE_CURRENT_SOURCE_DIR})

pkg_check_modules(VIPS REQUIRED vips-cpp)
if (VIPS_FOUND)
    message (STATUS "VIPS_INCLUDE_DIRS=${VIPS_INCLUDE_DIRS}")
    message (STATUS "VIPS_LIBRARIES=${VIPS_LIBRARIES}")
    message (STATUS "VIPS_LIBRARY_DIRS=${VIPS_LIBRARY_DIRS}")
    message (STATUS "VIPS_LDFLAGS=${VIPS_LDFLAGS}")
    link_directories(${VIPS_LIBRARY_DIRS})
    include_directories(${VIPS_INCLUDE_DIRS})
endif()

pkg_check_modules(GLIB2 REQUIRED glib-2.0)
if (GLIB2_FOUND)
    message (STATUS "GLIB2_INCLUDE_DIRS=${GLIB2_INCLUDE_DIRS}")
    message (STATUS "GLIB2_LIBRARIES=${GLIB2_LIBRARIES}")
    message (STATUS "GLIB2_LIBRARY_DIRS=${GLIB2_LIBRARY_DIRS}")
    message (STATUS "GLIB2_LDFLAGS=${GLIB2_LDFLAGS}")
    link_directories(${GLIB2_LIBRARY_DIRS})
    include_directories(${GLIB2_INCLUDE_DIRS})
endif()

# source files
FILE(GLOB Sources ./src/*.h)
FILE(GLOB Sources ./src/*.cpp)

# executable
add_executable(glitchfuck ${Sources})
target_link_libraries(glitchfuck ${GLIB2_LDFLAGS} ${VIPS_LDFLAGS})