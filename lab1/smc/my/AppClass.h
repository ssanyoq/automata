#ifndef _H_THECONTEXT
#define _H_THECONTEXT

#include "AppClass_sm.h"
#include <map>
#include <fstream>

#ifdef CRTP
class AppClass : public AppClassContext<AppClass>
#else
class AppClass
#endif
{
private:
#ifndef CRTP
    AppClassContext _fsm;
#endif

    bool isAcceptable;

	std::map<std::string, int> servers;
public:
    AppClass();
        // Default constructor.

    ~AppClass() {};
        // Destructor.

    bool CheckString(std::string);
        // Checks if the string is acceptable.

    void PrintMap(std::ostream &);

    inline void Acceptable()
    { isAcceptable = true; };

    inline void Unacceptable()
    { isAcceptable = false; };
        // State map actions.
};

#endif
