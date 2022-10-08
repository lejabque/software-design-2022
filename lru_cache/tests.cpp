#include "gtest/gtest.h"

#include "lru_cache.h"

namespace {
struct test_object {
  explicit test_object(int val) : value(val) {}
  int value;
};

} // namespace

TEST(lru_tests, simple_put_and_get) {
  LRU::LRUCache<int, test_object> cache(2);
  cache.Put(1, test_object(1));
  ASSERT_EQ(cache.Size(), 1);
  ASSERT_EQ(cache.Get(1).value, 1);
  ASSERT_THROW(cache.Get(2), std::out_of_range);
}

TEST(lru_tests, put_overwrite) {
  LRU::LRUCache<int, test_object> cache(1);
  cache.Put(1, test_object(1));
  ASSERT_EQ(cache.Size(), 1);
  ASSERT_EQ(cache.Get(1).value, 1);

  cache.Put(1, test_object(2));
  ASSERT_EQ(cache.Size(), 1);
  ASSERT_EQ(cache.Get(1).value, 2);
}

TEST(lru_tests, put_replaces_outdated) {
  LRU::LRUCache<int, test_object> cache(2);
  cache.Put(1, test_object(1));
  cache.Put(2, test_object(2));
  cache.Put(3, test_object(3));

  // 1 is removed from cache
  ASSERT_EQ(cache.Size(), 2);
  ASSERT_THROW(cache.Get(1), std::out_of_range);

  ASSERT_EQ(cache.Get(3).value, 3);
  ASSERT_EQ(cache.Get(2).value, 2);

  cache.Put(4, test_object(4));

  // 3 is removed from cache
  ASSERT_EQ(cache.Size(), 2);
  ASSERT_THROW(cache.Get(3), std::out_of_range);
  ASSERT_EQ(cache.Get(2).value, 2)  ;
  ASSERT_EQ(cache.Get(4).value, 4);
}
