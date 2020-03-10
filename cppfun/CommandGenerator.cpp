//
// Created by selva on 3/10/20.
//
#include "CommandGenerator.h"

CommandFactory::CommandFactory(){
    CommandGeneratorMap[PlayerCommandID::SetParam] = boost::bind(&CommandFactory::GenerateInit,this,_1);
    CommandGeneratorMap[PlayerCommandID::Play] = boost::bind(&CommandFactory::GeneratePlay,this,_1);
    CommandGeneratorMap[PlayerCommandID::GetTime] = boost::bind(&CommandFactory::GenerateGetTime,this,_1);
    CommandGeneratorMap[PlayerCommandID::Killservice] = boost::bind(&CommandFactory::GenerateKill,this,_1);
    CommandGeneratorMap[PlayerCommandID::Reset] = boost::bind(&CommandFactory::GenerateReset,this,_1);
    CommandGeneratorMap[PlayerCommandID::BCSH] = boost::bind(&CommandFactory::GenerateBCSH,this,_1);
    CommandGeneratorMap[PlayerCommandID::Stop] = boost::bind(&CommandFactory::GenerateStop,this,_1);
    CommandGeneratorMap[PlayerCommandID::IFwd] = boost::bind(&CommandFactory::GenerateIfwd,this,_1);
    CommandGeneratorMap[PlayerCommandID::IRwd] = boost::bind(&CommandFactory::GenerateIrwd,this,_1);
    CommandGeneratorMap[PlayerCommandID::Pause] = boost::bind(&CommandFactory::GeneratePause,this,_1);
    CommandGeneratorMap[PlayerCommandID::Resume] = boost::bind(&CommandFactory::GenerateResume,this,_1);
    CommandGeneratorMap[PlayerCommandID::Seek] = boost::bind(&CommandFactory::GenerateSeek,this,_1);
}
PlayerCommand CommandFactory::makeCommand(PlayerCommandID cmd_id, Iris2Msg msg) {
    auto itr = CommandGeneratorMap.find(cmd_id);
    if(itr != CommandGeneratorMap.end()){
        return itr->second(msg);
    }
    return GenerateNone(msg);
}
PlayerCommand CommandFactory::GenerateInit(Iris2Msg msg) {
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::SetParam;
    std::cout << "Generating Init Command";
    return cmd;
}
PlayerCommand CommandFactory::GeneratePlay(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Play;
    std::cout << "Generating Play Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateKill(Iris2Msg msg) {
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Killservice;
    std::cout << "Generating kill Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateReset(Iris2Msg msg) {
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Reset;
    std::cout << "Generating Reset Command";
    return cmd;

}
PlayerCommand CommandFactory::GenerateBCSH(Iris2Msg msg) {
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::BCSH;
    std::cout << "Generating BCSH Command";
    return cmd;

}
PlayerCommand CommandFactory::GenerateSeek(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Seek;
    std::cout << "Generating Seek Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateStop(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Stop;
    std::cout << "Generating Stop Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateResume(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Resume;
    std::cout << "Generating Resume Command";
    return cmd;
}
PlayerCommand CommandFactory::GeneratePause(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::Pause;
    std::cout << "Generating Pause Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateIfwd(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::IFwd;
    std::cout << "Generating IFwd Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateIrwd(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::IRwd;
    std::cout << "Generating IRwd Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateNone(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::None;
    std::cout << "Generating None Command";
    return cmd;
}
PlayerCommand CommandFactory::GenerateGetTime(Iris2Msg msg){
    InitCommand cmd;
    cmd.commandID = PlayerCommandID::GetTime;
    std::cout << "Generating Get Time Command";
    return cmd;
}