#include <iostream>
#include <string>
#include <fstream>
#include <chrono>
#include "AppClass.h"

using namespace statemap;

int main(int argc, char *argv[])
{
    AppClass thisContext;
    int retcode = 0;

    if (argc < 2)
    {
        std::cerr << "No string to check." << std::endl;
        retcode = 2;
    }
    else if (argc > 2)
    {
        std::cerr << "Only one argument is accepted." << std::endl;
        retcode = 3;
    }
    else
    {
    	std::ifstream fdata;
    	fdata.open(argv[1]);
    	if (!fdata)
    		throw std::runtime_error("Failed to open the file.");
    	auto start{std::chrono::steady_clock::now()};
    	while (!fdata.eof())
    	{
    		std::string str;
    		fdata >> str;
	        try
	        {
	            thisContext.CheckString(str);
	        }
	        catch (const SmcException &smcex)
	        {
	            std::cout << "not acceptable - "
	                 << smcex.what()
	                 << '.'
	                 << std::endl;

	            retcode = 1;
	        }    		
    	}
		auto end{std::chrono::steady_clock::now()};
		std::ofstream res;
		res.open("smc_res.txt");
		if (!res) {std::cout << "Error! Failed to open smc_res.txt." << std::endl; return 2;}
    	thisContext.PrintMap(res);
    	std::chrono::duration<double> elapsed_seconds{end - start};
    	std::cout << "Smc - Elapsed time: " << elapsed_seconds.count() << std::endl;

    }

    return retcode;
}
