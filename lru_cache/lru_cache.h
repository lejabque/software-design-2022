#include <cassert>
#include <list>
#include <unordered_map>

namespace LRU {
namespace details {

template <typename F>
struct DeferredCall {
  explicit DeferredCall(F f) : f_(f) {}
  ~DeferredCall() {
    f_();
  }
  F f_;
};

} // namespace details
template <typename K, typename V>
class LRUCache {
public:
  explicit LRUCache(size_t capacity) : Capacity_(capacity) {}

  const V& Get(const K& key) {
    auto it = Cache_.find(key);
    if (it == Cache_.end()) {
      throw std::out_of_range("No such key");
    }
    MoveToFront(it->second.ListIt);

    CheckSizeInv();
    assert(RecentlyUsed_.front() == key);
    return it->second.Value;
  }

  void Put(K key, V value) {
    auto [it, inserted] = Cache_.try_emplace(key, std::move(value));
    if (inserted) {
      RecentlyUsed_.push_front(std::move(key));
      it->second.ListIt = RecentlyUsed_.begin();
    } else {
      it->second.Value = std::move(value);
      MoveToFront(it->second.ListIt);
    }

    if (RecentlyUsed_.size() > Capacity_) {
      auto mostRecently = std::prev(RecentlyUsed_.end());
      Cache_.erase(*mostRecently);
      RecentlyUsed_.erase(mostRecently);
    }

    CheckSizeInv();
    assert(RecentlyUsed_.front() == key);
  }

  size_t Size() const {
    CheckSizeInv();
    return Cache_.size();
  }

private:
  using ListIter = typename std::list<K>::const_iterator;

  void CheckSizeInv() const {
    assert(0 <= Cache_.size() && Cache_.size() <= Capacity_);
    assert(RecentlyUsed_.size() == Cache_.size());
  }

  void MoveToFront(ListIter it) {
#ifndef NDEBUG
    auto listBefore = RecentlyUsed_;
#endif

    RecentlyUsed_.splice(RecentlyUsed_.begin(), RecentlyUsed_, it);

#ifndef NDEBUG
    assert(RecentlyUsed_.front() == *it); // moved to front
    // other part of list not changed
    listBefore.remove(*it);
    assert(std::equal(listBefore.begin(), listBefore.end(),
                      std::next(RecentlyUsed_.begin())));
#endif
  }

  struct Node {
    explicit Node(V val) : Value(std::move(val)) {}

    V Value;
    ListIter ListIt; // iterator to node in recently used list
  };

  std::unordered_map<K, Node> Cache_;
  // we copy key here, but it can be implemented more efficiently, e.g. using
  // intrusive list
  std::list<K> RecentlyUsed_;
  size_t Capacity_;
};
} // namespace LRU
