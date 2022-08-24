#include "random.hpp"
#include <chrono>

using namespace std::chrono;

int nanosecondSeed() {
  high_resolution_clock::time_point now = high_resolution_clock::now();
  auto now_ns = time_point_cast<nanoseconds>(now);
  auto epoch = now_ns.time_since_epoch();
  auto seed = duration_cast<nanoseconds>(epoch);
  return seed.count();
}