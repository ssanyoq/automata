%option c++
%option noyywrap
%{
	#include <map>
	#include <iostream>
	#include <fstream>
	#include <string>
	#include <chrono>
	std::map<std::string, int> servers;
    int size = 0;
    std::string server;
%}
jid [a-zA-Z0-9]+@
servern [a-zA-Z0-9]+
domainname \.[a-zA-Z0-9]{1,5}
%x USERNAME SERVERNAME ADDRESS EOL
%%
<INITIAL>{
	gtalk:(talk|chat)\?jid= {BEGIN(USERNAME);}
	. ;
	\n {return 0;}
}
<USERNAME>{
    ({jid}) {
        BEGIN(SERVERNAME);
    }
    . ;
    \n {return 0;}
}

<SERVERNAME>{
    ({servern}) {
        server = yytext;
        BEGIN(ADDRESS);
    }
    . ;
    \n {return 0;}
}
<ADDRESS>{
    ({domainname}) {
        if (servers.count(yytext) != 0) servers[yytext] += 1;
		else servers[yytext] = 1;
        BEGIN(EOL);
    }
    . ;
    \n {return 0;}
}
<EOL>{
    \n {
      BEGIN(INITIAL);
    }
    . ;
}

%%

int main()
{
	std::ifstream fdata("input.txt");
	if (!fdata)
	{
		std::cout << "Failed to open the file" << std::endl;
		return 2;
	}
	yyFlexLexer lex(fdata, std::cout);
	const auto start{std::chrono::steady_clock::now()};
	while(lex.yylex());
	const auto end {std::chrono::steady_clock::now()};
	fdata.close();
	std::ofstream res;
	res.open("flex_res.txt");
	for (auto it = servers.begin(); it != servers.end(); ++it)
		res << it->first << ": " << it->second << '\n';
	res.close();
	const std::chrono::duration<double> elapsed_seconds {end-start};
	std::cout << "Flex - Elapsed time: " << elapsed_seconds.count() << std::endl;
}

