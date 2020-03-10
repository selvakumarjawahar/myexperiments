#include <iostream>
#include "CommandGenerator.h"
#include <boost/variant.hpp>

int main() {
    CommandFactory factory;
    Iris2Msg msg;
    auto cmd = factory.makeCommand(PlayerCommandID::SetParam,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Play,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Stop,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::IFwd,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::IRwd,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Seek,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Pause,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Resume,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::GetTime,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Killservice,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::Reset,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::BCSH,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';
    cmd = factory.makeCommand(PlayerCommandID::None,msg);
    std::cout << '\n' << "Command ID = " << (int)boost::apply_visitor(GetCommandID(),cmd) << '\n';

    return 0;
}
