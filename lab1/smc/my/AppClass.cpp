#ifdef _MSC_VER
#pragma warning(disable: 4355)
#endif

#include "AppClass.h"
#include <cstring>
#include <iostream>
#include <sstream>

AppClass::AppClass()
#ifdef CRTP
: isAcceptable(false)
#else
: _fsm(*this),
  isAcceptable(false)
#endif
{
#ifdef FSM_DEBUG
#ifdef CRTP
    setDebugFlag(true);
#else
    _fsm.setDebugFlag(true);
#endif
#endif
}

bool AppClass::CheckString(std::string theString)
{
	const std::string num = "0123456789";
	const std::string alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

	int domainSize = 0;
	bool isPeriod = false;

	std::string gtalk, meetType, jid;
	std::copy(theString.begin(), theString.begin() + 6, std::back_inserter(gtalk));
	std::copy(theString.begin() + 6, theString.begin() + 10, std::back_inserter(meetType));
	std::copy(theString.begin() + 10, theString.begin() + 15, std::back_inserter(jid));
#ifdef CRTP
    enterStartState();
    Reset();

    if (gtalk.compare("gtalk:") == 0 && (meetType.compare("chat") == 0 || meetType.compare("talk") == 0) && jid.compare("?jid=") == 0) {
    	Prefix();
    }
    else
    	Unknown();
    for (unsigned long i = 15; i < theString.size(); i++)
    {
        if (num.find(theString[i]) != std::string::npos) // numbers
        {
        	Num();
        } else if (alpha.find(theString[i]) != std::string::npos) // alpha
    	{
        	if (isPeriod) {
          		domainSize++;
	            if (domainSize > 5) {
            		AlphaLimit();
	                continue;
	            }
            }
            Alpha();
    	} else {
          switch (theString[i]) {
            case '.':
              Period();
              isPeriod = true;
              break;
            case '@':
              At();
              break;
            default:
              Unknown();
          }
        }
    }

    EOS();
#else
	_fsm.enterStartState();
	_fsm.Reset();
    if (gtalk.compare("gtalk:") == 0 && (meetType.compare("chat") == 0 || meetType.compare("talk") == 0) && jid.compare("?jid=") == 0) {
		_fsm.Prefix();
	}
	else
		_fsm.Unknown();
	for (unsigned long i = 15; i < theString.size(); i++)
	{
		if (num.find(theString[i]) != std::string::npos) // numbers
		{
			_fsm.Num();
		} else if (alpha.find(theString[i]) != std::string::npos) // alpha
		{
			if (isPeriod) {
				domainSize++;
				if (domainSize > 5) {
					_fsm.AlphaLimit();
					continue;
				}
			}
			_fsm.Alpha();
		} else {
			switch (theString[i]) {
				case '.':
					_fsm.Period();
					isPeriod = true;
				break;
				case '@':
					_fsm.At();
				break;
				default:
					_fsm.Unknown();
			}
		}
	}

	_fsm.EOS();
#endif
	if (isAcceptable)
	{
		size_t start = theString.find("@");
		size_t end = theString.find(".");
		std::string serverName;
		std::copy(theString.begin() + start + 1, theString.begin() + end, std::back_inserter(serverName));
		std::istringstream nameStream(serverName);
		std::string token;
		while (std::getline(nameStream, token, ','))
		{
			if (servers.count(token) != 0)
				servers[token] += 1;
			else
				servers[token] = 1;
		}
	}
    return isAcceptable;
}

void AppClass::PrintMap(std::ostream &c)
{
	for (auto it = servers.begin(); it != servers.end(); it++)
	{
		c << it->first << ": " << it->second << std::endl;
	}
}
