//
// Created by selva on 3/21/21.
//

#include <iostream>
#include <memory>
#include <thread>
#include <chrono>
#include <future>

struct sleeper {
  int class_id;
  int sleep_sec;
  sleeper(int id, int sleep_sec):
                                   class_id(id),
                                   sleep_sec(sleep_sec){}
  void go_sleep(){
    std::cout << "sleeper ID = " << class_id << " going to sleep for "
        <<sleep_sec<<" seconds"<<'\n';
    std::this_thread::sleep_for(std::chrono::seconds(sleep_sec));
  }
  ~sleeper(){
    std::cout << "sleeper ID = " << class_id << " going to destroy now" << '\n';
  }
};

int call_sleeper(int id, int sleep) {
  auto ptr = std::make_shared<sleeper>(id,sleep);
  ptr ->go_sleep();
  return sleep;
}

int main() {
  for(int i = 0; i<30 ; i++){
    auto ret = std::async(std::launch::async ,[i](){return call_sleeper(i,i+1);});
  }
  auto myfuture = std::async([](){return call_sleeper(30,32);});
  myfuture.get();
}