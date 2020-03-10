//
// Created by selva on 3/10/20.
//

#ifndef CPPFUN_COMMANDGENERATOR_H
#define CPPFUN_COMMANDGENERATOR_H

#include <string>
#include <iostream>
#include <map>

#include <boost/variant.hpp>
#include <boost/function.hpp>
#include <boost/bind.hpp>

enum class PlayerCommandID {
    SetParam,
    Play,
    Stop,
    IFwd,
    IRwd,
    Seek,
    Pause,
    Resume,
    GetTime,
    GetEOF,
    Killservice,
    Reset,
    BCSH,
    None
};

enum class VIDEO_CODEC : char {
    H264,
    H265,
    MP4
};

enum class AUDIO_CODEC : char {
    AC3,
    AAC,
    PCM,
    AC3PASSTHRU
};
enum class RESOLUTION : char {
    RES1080p, RES720p, RES1080i
};

enum class FRAME_RATE : char { //! As more frame rates are added this enum
    //! should grow
    FPS24 = 24,
    FPS25 = 25,
    FPS30 = 30,
    FPS50 = 50,
    FPS60 = 60
};

//CRTP Design Pattern to provide Interface to interact with commands
template <typename T>
class CommonCommandInterface {
public:
    PlayerCommandID getCommandID(){
        T& command_derived = static_cast<T&>(*this);
        return command_derived.commandID;
    }
private:
    CommonCommandInterface(){};
    friend T;
};

struct Preview {
    bool previewMode;
    uint32_t textColor;
    uint32_t width;
    uint32_t height;
    uint32_t xPos;
    uint32_t yPos;
    bool border;
    uint32_t borderColor;
    std::string fontPath;
    uint32_t fontSize;
    std::string previewText;
};

struct RunningMark {
    bool mode;
    uint32_t color;
    uint32_t width;
    uint32_t height;
    bool border;
    uint32_t borderColor;
};

struct GridMark {
    bool mode;
    uint32_t markColor;
    uint32_t spaceColor;
    uint32_t width;
    uint32_t height;
    uint32_t rows;
    uint32_t cols;
    bool border;
    uint32_t borderColor;
};

struct InitCommand: public CommonCommandInterface<InitCommand> {
    PlayerCommandID commandID;
    bool videoPresent;
    bool audioPresent;
    uint64_t projectorSetupTime;
    VIDEO_CODEC videoCodec;
    AUDIO_CODEC audioCodec;
    RESOLUTION resolution;
    FRAME_RATE frameRate;
    bool setAudioDecoder;
    bool streamReader;
    int theatreID;
    Preview previewInfo;
    RunningMark rmInfo;
    GridMark gmInfo;
};

struct PlayCommand: public CommonCommandInterface<PlayCommand> {
    PlayerCommandID commandID;
    std::string videoPath;
    std::string audioPath;
    bool encrypted;
    std::string key;
    uint32_t seekToFrame;
    bool doAVPlayback;
    RESOLUTION resolution;
    FRAME_RATE frameRate;
    uint32_t mediaID;
};

struct NullCommand: public CommonCommandInterface<NullCommand> {
    PlayerCommandID commandID;
};

struct KillCommand: public CommonCommandInterface<KillCommand> {
    PlayerCommandID commandID;
};

struct ResetCommand: public CommonCommandInterface<ResetCommand> {
    PlayerCommandID commandID;
};

struct BCSHCommand: public CommonCommandInterface<BCSHCommand> {
    PlayerCommandID commandID;
    uint8_t id;
    uint32_t value;
    bool reset;
};

struct TrickplayCommand: public CommonCommandInterface<TrickplayCommand> {
    PlayerCommandID commandID;
    uint32_t speed;
    uint32_t playbackState;
    uint32_t maxTrickDuration;
    uint32_t seekOpts;
};

struct Iris2Msg {
    std::string message;
};

using PlayerCommand = boost::variant<InitCommand,
        PlayCommand,
        NullCommand,
        KillCommand,
        ResetCommand,
        BCSHCommand,
        TrickplayCommand>;

class CommandFactory {
public:
    CommandFactory();
    PlayerCommand makeCommand(PlayerCommandID cmd_id,Iris2Msg msg );
private:
    std::map<PlayerCommandID,boost::function<PlayerCommand(Iris2Msg msg)>> CommandGeneratorMap;
    PlayerCommand GenerateInit(Iris2Msg msg);
    PlayerCommand GeneratePlay(Iris2Msg msg);
    PlayerCommand GenerateKill(Iris2Msg msg);
    PlayerCommand GenerateReset(Iris2Msg msg);
    PlayerCommand GenerateBCSH(Iris2Msg msg);
    PlayerCommand GeneratePause(Iris2Msg msg);
    PlayerCommand GenerateResume(Iris2Msg msg);
    PlayerCommand GenerateStop(Iris2Msg msg);
    PlayerCommand GenerateSeek(Iris2Msg msg);
    PlayerCommand GenerateIfwd(Iris2Msg msg);
    PlayerCommand GenerateIrwd(Iris2Msg msg);
    PlayerCommand GenerateGetTime(Iris2Msg msg);
    PlayerCommand GenerateNone(Iris2Msg msg);

};

class GetCommandID
        : public boost::static_visitor<PlayerCommandID>
{
public:

    template <typename T>
    PlayerCommandID operator()( T & cmd ) const
    {
        return cmd.getCommandID();
    }

};

#endif //CPPFUN_COMMANDGENERATOR_H
