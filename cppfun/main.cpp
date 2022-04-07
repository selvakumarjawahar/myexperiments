#include <iostream>
#include <vector>
#include <array>
#include <algorithm>

using namespace std;

/*
class Solution {
public:
  bool checkValid(vector<vector<int>>& matrix) {
    for (const auto& row : matrix )
    {
      std::vector<int> num_map(matrix.size()+1, 0);
      num_map[0] = 1;
      for (const auto& elem : row)
      {
        num_map[elem] = 1;
      }
      if(!std::all_of(num_map.begin(),num_map.end(),
                       [](int i){return i == 1;}))
      {
        return false;
      }
    }

    for(int j = 0 ;j < matrix.size(); j++)
    {
      std::vector<int> num_map(matrix.size()+1, 0);
      num_map[0] = 1;
      for(int i = 0; i< matrix.size(); i++)
      {
        num_map[matrix[i][j]] = 1;
      }
      if(!std::all_of(num_map.begin(),num_map.end(),
                       [](int i){return i == 1;}))
      {
        return false;
      }
    }
    return true;
  }
};
*/

class Solution {
public:
  int minSwaps(vector<int>& nums) {
    vector<int> gaps;
    int state = 0;
    int gap_count = 0;
    for(const auto& elem : nums)
    {
      switch(state)
      {
      case 0: // find one
        if(elem == 1)
          state = 1;
        break;
      case 1: // find 0
        if(elem == 0)
        {
          state = 2;
          gap_count++;
        }
        break;
      case 2: // find 1
        if(elem == 1)
        {
          gaps.push_back(gap_count);
          state = 1;
          gap_count = 0;
        }
        else
        {
          gap_count++;
        }
        break;
      }
    }
    auto itr1 = find(nums.begin(),nums.end(),1);
    if(itr1 != nums.end()) {
      auto itr2 = find(nums.rbegin(), nums.rend(), 1);
      if (itr2.base() != nums.begin())
      {
        auto l1 = distance(itr2.base(),nums.end());
        auto l2 = distance(nums.begin(),itr1);
        gaps.push_back(abs(l2-l1));
      }

    }
    return *min_element(gaps.begin(),gaps.end());
  }
};

int main()
{
  Solution solution;
  vector<int> tc1{0,1,0,1,1,0,0};
  cout << "tc1 - " << solution.minSwaps(tc1) << '\n';
  vector<int> tc2{0,1,1,1,0,0,1,1,0};
  cout << "tc2 - " << solution.minSwaps(tc2) << '\n';
  vector<int> tc3{1,1,0,0,1};
  cout << "tc3 - " << solution.minSwaps(tc3) << '\n';

  return 0;
}