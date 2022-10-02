class LRUCache {
public:
  LRUCache(int capacity) : Capacity_(capacity) {}

  int get(int key) {
    auto it = Cache_.find(key);
    if (it == Cache_.end()) {
      throw std::out_of_range("Key not found");
    }
    MoveToFront(it->second.ListIt);
    return it->second.Value;
  }

  void put(int key, int value) {
    auto [it, inserted] = Cache_.try_emplace(key, std::move(value));
    if (inserted) {
      Lru_.push_front(std::move(key));
      it->second.ListIt = Lru_.begin();
    } else {
      it->second.Value = std::move(value);
      MoveToFront(it->second.ListIt);
    }
    if (Lru_.size() > Capacity_) {
      EraseOutdated();
    }
  }

private:
  using K = int;
  using V = int;

private:
  using ListIter = typename std::list<K>::const_iterator;

  void MoveToFront(ListIter it) {
    Lru_.splice(Lru_.begin(), Lru_, it);
  }

  void EraseOutdated() {
    auto it = std::prev(Lru_.end());
    Cache_.erase(*it);
    Lru_.erase(it);
  }

  struct Node {
    explicit Node(V val) : Value(std::move(val)) {}

    V Value;
    ListIter ListIt;
  };

  // we copy key here but it can be implemented more ef
  std::unordered_map<K, Node> Cache_;
  std::list<K> Lru_;
  int Capacity_;
};

/**
 * Your LRUCache object will be instantiated and called as such:
 * LRUCache* obj = new LRUCache(capacity);
 * int param_1 = obj->get(key);
 * obj->put(key,value);
 */