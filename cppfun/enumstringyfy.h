/*
 * The following code is taken from the stackoverflow answer
 * https://stackoverflow.com/questions/5093460/how-to-convert-an-enum-type-variable-to-a-string
 * The Code has been modified to take enum class instead of simple enum
 * This code uses Boost Preprocessor Library for generation of switch case statements and to stringify
 */

#include <boost/preprocessor.hpp>

#define X_DEFINE_ENUM_WITH_STRING_CONVERSIONS_TOSTRING_CASE(r,enumtype, elem)    \
    case enumtype::elem : return BOOST_PP_STRINGIZE(elem);

#define DEFINE_ENUM_WITH_STRING_CONVERSIONS(enumtype, enumerators)                \
    enum class enumtype {                                                               \
        BOOST_PP_SEQ_ENUM(enumerators)                                        \
    };                                                                        \
                                                                              \
    inline const char* ToString(enumtype v)                                       \
    {                                                                         \
        switch (v)                                                            \
        {                                                                     \
            BOOST_PP_SEQ_FOR_EACH(                                            \
                X_DEFINE_ENUM_WITH_STRING_CONVERSIONS_TOSTRING_CASE,          \
                enumtype,                                                         \
                enumerators                                                   \
            )                                                                 \
            default: return "[Unknown " BOOST_PP_STRINGIZE(name) "]";         \
        }                                                                     \
    }