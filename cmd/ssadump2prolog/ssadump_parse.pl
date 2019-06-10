%:- set_prolog_flag(double_quotes, chars).

:- use_module(library(readutil)).
:- use_module(library(dcg/basics)).

%%%%%%
% ----------------
prologue -->
	comment("Name", Name),
	comment("Package", Pkg),
	comment("Location", Loc),
	% ---
	sig(Rcv, As, Rs), maplist(is_list, [Rcv, As, Rs]).
	% ---



% ----------------
comment(Lbl, Txt) -->
	"#", blanks, string(Lbl),
	":", blanks, string_without("\n", Txt).

sig(FuncName, (RcvName, RcvType), As, Rs) -->
	"func", "(", 
		string_without(" ", RcvName),
		string_without(" ", RcvType),
	")",
	string_without(" ", FuncName), 
	"(",	sig_args(As), ")",
	"(",	sig_rets(Rs), ")".

block_label(Id, PredsNum, SuccsNum, Comment) --> 
	digits(Id), ":", blanks,
	%%%
	string_without(" ", Comment), blanks, 
	"P:", digits(PredsNum), blanks,
	"S:", digits(SuccsNum).

block_body --> instrs.

instrs --> stmt(S).
instrs --> stmt(S), instrs.
instrs --> expr(E, Ty), instrs.

%%% TODO: more intelligent var-type detector %%%%%%%%%%%%%%%%
/***\
||51||   [*func(Request, int) (bool, time.Duration)]
|	|	
|52||    [func(Request, int) (bool, time.Duration)]
\***/
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%


%%%%%
main(File) :- 
	read_file_to_codes(File, CC0, []),
	string_codes(SS, CC0),
	phrase(comment(Lbl, Txt), CC0, Unparsed),
	writeln(SS),
	%phrase_from_file(comment(Lbl, Txt), File),
	write("\n\n\n--------------\n\n\n"),

	format("LBL: ~s ;; TXT: ~s ;;\nUNP: ~s\n\n----\n",
			[Lbl, Txt, Unparsed]).

main :- main('C:/Users/XPinguin/_projects/_learning/prolog/testdata/sample_ssa.txt').

