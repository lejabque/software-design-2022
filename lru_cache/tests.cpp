#include "gtest/gtest.h"

#include "lru_cache.h"

int sum(int a, int b) {
  // In normal tasks, never change files with tests but here you can :)
  return a + b; // TODO: fix me
}

TEST(something, simple) {
  ASSERT_EQ(2, sum(1, 1));
}
