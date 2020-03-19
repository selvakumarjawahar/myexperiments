#include <iostream>
#include <torch/torch.h>

int main(){
    auto size = torch::IntArrayRef{3,3};
    auto topts = torch::TensorOptions().dtype(torch::kFloat32).device(torch::kCPU).layout(torch::kStrided).requires_grad(
            false);
    auto x = torch::empty(size,topts);
    auto y = torch::rand(size,topts);
    auto z = torch::ones(size,topts);
    auto a = torch::tensor({{4,5,6},{6,5,4},{2,3,4}},topts);
    auto result = torch::empty(size,topts);
    torch::add_out(result,z,a);
    auto res_acc = result.accessor<float,2>();
    std::cout << a.reshape({1,9}) << '\n';
    if (torch::cuda::is_available()) {
        std::cout << "CUDA is available" << std::endl;
    } else{
        std::cout << "CUDA not available" << std::endl;
    }

}
