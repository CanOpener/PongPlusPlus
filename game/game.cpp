#include <iostream>
#include <mutex>
#include <deque>
using namespace std;

typedef unsigned char BYTE

// thread safe queuing
mutex mu;
deque<BYTE*> callQueue;

void add(BYTE* bytes) {
    lock_guard<mutex> locker(mu);
    callQueue.push_back(bytes);
}
deque<BYTE*> getAll() {
    lock_guard<mutex> locker(mu);
    deque<BYTE*> ret = callQueue;
    callQueue.clear();
    return ret;
}


int main() {
    cout << "Hello World!" << endl;
}
