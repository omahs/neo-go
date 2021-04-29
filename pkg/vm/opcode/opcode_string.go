// Code generated by "stringer -type=Opcode -linecomment"; DO NOT EDIT.

package opcode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PUSHINT8-0]
	_ = x[PUSHINT16-1]
	_ = x[PUSHINT32-2]
	_ = x[PUSHINT64-3]
	_ = x[PUSHINT128-4]
	_ = x[PUSHINT256-5]
	_ = x[PUSHA-10]
	_ = x[PUSHNULL-11]
	_ = x[PUSHDATA1-12]
	_ = x[PUSHDATA2-13]
	_ = x[PUSHDATA4-14]
	_ = x[PUSHM1-15]
	_ = x[PUSH0-16]
	_ = x[PUSHF-16]
	_ = x[PUSH1-17]
	_ = x[PUSHT-17]
	_ = x[PUSH2-18]
	_ = x[PUSH3-19]
	_ = x[PUSH4-20]
	_ = x[PUSH5-21]
	_ = x[PUSH6-22]
	_ = x[PUSH7-23]
	_ = x[PUSH8-24]
	_ = x[PUSH9-25]
	_ = x[PUSH10-26]
	_ = x[PUSH11-27]
	_ = x[PUSH12-28]
	_ = x[PUSH13-29]
	_ = x[PUSH14-30]
	_ = x[PUSH15-31]
	_ = x[PUSH16-32]
	_ = x[NOP-33]
	_ = x[JMP-34]
	_ = x[JMPL-35]
	_ = x[JMPIF-36]
	_ = x[JMPIFL-37]
	_ = x[JMPIFNOT-38]
	_ = x[JMPIFNOTL-39]
	_ = x[JMPEQ-40]
	_ = x[JMPEQL-41]
	_ = x[JMPNE-42]
	_ = x[JMPNEL-43]
	_ = x[JMPGT-44]
	_ = x[JMPGTL-45]
	_ = x[JMPGE-46]
	_ = x[JMPGEL-47]
	_ = x[JMPLT-48]
	_ = x[JMPLTL-49]
	_ = x[JMPLE-50]
	_ = x[JMPLEL-51]
	_ = x[CALL-52]
	_ = x[CALLL-53]
	_ = x[CALLA-54]
	_ = x[CALLT-55]
	_ = x[ABORT-56]
	_ = x[ASSERT-57]
	_ = x[THROW-58]
	_ = x[TRY-59]
	_ = x[TRYL-60]
	_ = x[ENDTRY-61]
	_ = x[ENDTRYL-62]
	_ = x[ENDFINALLY-63]
	_ = x[RET-64]
	_ = x[SYSCALL-65]
	_ = x[DEPTH-67]
	_ = x[DROP-69]
	_ = x[NIP-70]
	_ = x[XDROP-72]
	_ = x[CLEAR-73]
	_ = x[DUP-74]
	_ = x[OVER-75]
	_ = x[PICK-77]
	_ = x[TUCK-78]
	_ = x[SWAP-80]
	_ = x[ROT-81]
	_ = x[ROLL-82]
	_ = x[REVERSE3-83]
	_ = x[REVERSE4-84]
	_ = x[REVERSEN-85]
	_ = x[INITSSLOT-86]
	_ = x[INITSLOT-87]
	_ = x[LDSFLD0-88]
	_ = x[LDSFLD1-89]
	_ = x[LDSFLD2-90]
	_ = x[LDSFLD3-91]
	_ = x[LDSFLD4-92]
	_ = x[LDSFLD5-93]
	_ = x[LDSFLD6-94]
	_ = x[LDSFLD-95]
	_ = x[STSFLD0-96]
	_ = x[STSFLD1-97]
	_ = x[STSFLD2-98]
	_ = x[STSFLD3-99]
	_ = x[STSFLD4-100]
	_ = x[STSFLD5-101]
	_ = x[STSFLD6-102]
	_ = x[STSFLD-103]
	_ = x[LDLOC0-104]
	_ = x[LDLOC1-105]
	_ = x[LDLOC2-106]
	_ = x[LDLOC3-107]
	_ = x[LDLOC4-108]
	_ = x[LDLOC5-109]
	_ = x[LDLOC6-110]
	_ = x[LDLOC-111]
	_ = x[STLOC0-112]
	_ = x[STLOC1-113]
	_ = x[STLOC2-114]
	_ = x[STLOC3-115]
	_ = x[STLOC4-116]
	_ = x[STLOC5-117]
	_ = x[STLOC6-118]
	_ = x[STLOC-119]
	_ = x[LDARG0-120]
	_ = x[LDARG1-121]
	_ = x[LDARG2-122]
	_ = x[LDARG3-123]
	_ = x[LDARG4-124]
	_ = x[LDARG5-125]
	_ = x[LDARG6-126]
	_ = x[LDARG-127]
	_ = x[STARG0-128]
	_ = x[STARG1-129]
	_ = x[STARG2-130]
	_ = x[STARG3-131]
	_ = x[STARG4-132]
	_ = x[STARG5-133]
	_ = x[STARG6-134]
	_ = x[STARG-135]
	_ = x[NEWBUFFER-136]
	_ = x[MEMCPY-137]
	_ = x[CAT-139]
	_ = x[SUBSTR-140]
	_ = x[LEFT-141]
	_ = x[RIGHT-142]
	_ = x[INVERT-144]
	_ = x[AND-145]
	_ = x[OR-146]
	_ = x[XOR-147]
	_ = x[EQUAL-151]
	_ = x[NOTEQUAL-152]
	_ = x[SIGN-153]
	_ = x[ABS-154]
	_ = x[NEGATE-155]
	_ = x[INC-156]
	_ = x[DEC-157]
	_ = x[ADD-158]
	_ = x[SUB-159]
	_ = x[MUL-160]
	_ = x[DIV-161]
	_ = x[MOD-162]
	_ = x[POW-163]
	_ = x[SQRT-164]
	_ = x[SHL-168]
	_ = x[SHR-169]
	_ = x[NOT-170]
	_ = x[BOOLAND-171]
	_ = x[BOOLOR-172]
	_ = x[NZ-177]
	_ = x[NUMEQUAL-179]
	_ = x[NUMNOTEQUAL-180]
	_ = x[LT-181]
	_ = x[LE-182]
	_ = x[GT-183]
	_ = x[GE-184]
	_ = x[MIN-185]
	_ = x[MAX-186]
	_ = x[WITHIN-187]
	_ = x[PACK-192]
	_ = x[UNPACK-193]
	_ = x[NEWARRAY0-194]
	_ = x[NEWARRAY-195]
	_ = x[NEWARRAYT-196]
	_ = x[NEWSTRUCT0-197]
	_ = x[NEWSTRUCT-198]
	_ = x[NEWMAP-200]
	_ = x[SIZE-202]
	_ = x[HASKEY-203]
	_ = x[KEYS-204]
	_ = x[VALUES-205]
	_ = x[PICKITEM-206]
	_ = x[APPEND-207]
	_ = x[SETITEM-208]
	_ = x[REVERSEITEMS-209]
	_ = x[REMOVE-210]
	_ = x[CLEARITEMS-211]
	_ = x[POPITEM-212]
	_ = x[ISNULL-216]
	_ = x[ISTYPE-217]
	_ = x[CONVERT-219]
}

const _Opcode_name = "PUSHINT8PUSHINT16PUSHINT32PUSHINT64PUSHINT128PUSHINT256PUSHAPUSHNULLPUSHDATA1PUSHDATA2PUSHDATA4PUSHM1PUSH0PUSH1PUSH2PUSH3PUSH4PUSH5PUSH6PUSH7PUSH8PUSH9PUSH10PUSH11PUSH12PUSH13PUSH14PUSH15PUSH16NOPJMPJMP_LJMPIFJMPIF_LJMPIFNOTJMPIFNOT_LJMPEQJMPEQ_LJMPNEJMPNE_LJMPGTJMPGT_LJMPGEJMPGE_LJMPLTJMPLT_LJMPLEJMPLE_LCALLCALL_LCALLACALLTABORTASSERTTHROWTRYTRY_LENDTRYENDTRY_LENDFINALLYRETSYSCALLDEPTHDROPNIPXDROPCLEARDUPOVERPICKTUCKSWAPROTROLLREVERSE3REVERSE4REVERSENINITSSLOTINITSLOTLDSFLD0LDSFLD1LDSFLD2LDSFLD3LDSFLD4LDSFLD5LDSFLD6LDSFLDSTSFLD0STSFLD1STSFLD2STSFLD3STSFLD4STSFLD5STSFLD6STSFLDLDLOC0LDLOC1LDLOC2LDLOC3LDLOC4LDLOC5LDLOC6LDLOCSTLOC0STLOC1STLOC2STLOC3STLOC4STLOC5STLOC6STLOCLDARG0LDARG1LDARG2LDARG3LDARG4LDARG5LDARG6LDARGSTARG0STARG1STARG2STARG3STARG4STARG5STARG6STARGNEWBUFFERMEMCPYCATSUBSTRLEFTRIGHTINVERTANDORXOREQUALNOTEQUALSIGNABSNEGATEINCDECADDSUBMULDIVMODPOWSQRTSHLSHRNOTBOOLANDBOOLORNZNUMEQUALNUMNOTEQUALLTLEGTGEMINMAXWITHINPACKUNPACKNEWARRAY0NEWARRAYNEWARRAY_TNEWSTRUCT0NEWSTRUCTNEWMAPSIZEHASKEYKEYSVALUESPICKITEMAPPENDSETITEMREVERSEITEMSREMOVECLEARITEMSPOPITEMISNULLISTYPECONVERT"

var _Opcode_map = map[Opcode]string{
	0:   _Opcode_name[0:8],
	1:   _Opcode_name[8:17],
	2:   _Opcode_name[17:26],
	3:   _Opcode_name[26:35],
	4:   _Opcode_name[35:45],
	5:   _Opcode_name[45:55],
	10:  _Opcode_name[55:60],
	11:  _Opcode_name[60:68],
	12:  _Opcode_name[68:77],
	13:  _Opcode_name[77:86],
	14:  _Opcode_name[86:95],
	15:  _Opcode_name[95:101],
	16:  _Opcode_name[101:106],
	17:  _Opcode_name[106:111],
	18:  _Opcode_name[111:116],
	19:  _Opcode_name[116:121],
	20:  _Opcode_name[121:126],
	21:  _Opcode_name[126:131],
	22:  _Opcode_name[131:136],
	23:  _Opcode_name[136:141],
	24:  _Opcode_name[141:146],
	25:  _Opcode_name[146:151],
	26:  _Opcode_name[151:157],
	27:  _Opcode_name[157:163],
	28:  _Opcode_name[163:169],
	29:  _Opcode_name[169:175],
	30:  _Opcode_name[175:181],
	31:  _Opcode_name[181:187],
	32:  _Opcode_name[187:193],
	33:  _Opcode_name[193:196],
	34:  _Opcode_name[196:199],
	35:  _Opcode_name[199:204],
	36:  _Opcode_name[204:209],
	37:  _Opcode_name[209:216],
	38:  _Opcode_name[216:224],
	39:  _Opcode_name[224:234],
	40:  _Opcode_name[234:239],
	41:  _Opcode_name[239:246],
	42:  _Opcode_name[246:251],
	43:  _Opcode_name[251:258],
	44:  _Opcode_name[258:263],
	45:  _Opcode_name[263:270],
	46:  _Opcode_name[270:275],
	47:  _Opcode_name[275:282],
	48:  _Opcode_name[282:287],
	49:  _Opcode_name[287:294],
	50:  _Opcode_name[294:299],
	51:  _Opcode_name[299:306],
	52:  _Opcode_name[306:310],
	53:  _Opcode_name[310:316],
	54:  _Opcode_name[316:321],
	55:  _Opcode_name[321:326],
	56:  _Opcode_name[326:331],
	57:  _Opcode_name[331:337],
	58:  _Opcode_name[337:342],
	59:  _Opcode_name[342:345],
	60:  _Opcode_name[345:350],
	61:  _Opcode_name[350:356],
	62:  _Opcode_name[356:364],
	63:  _Opcode_name[364:374],
	64:  _Opcode_name[374:377],
	65:  _Opcode_name[377:384],
	67:  _Opcode_name[384:389],
	69:  _Opcode_name[389:393],
	70:  _Opcode_name[393:396],
	72:  _Opcode_name[396:401],
	73:  _Opcode_name[401:406],
	74:  _Opcode_name[406:409],
	75:  _Opcode_name[409:413],
	77:  _Opcode_name[413:417],
	78:  _Opcode_name[417:421],
	80:  _Opcode_name[421:425],
	81:  _Opcode_name[425:428],
	82:  _Opcode_name[428:432],
	83:  _Opcode_name[432:440],
	84:  _Opcode_name[440:448],
	85:  _Opcode_name[448:456],
	86:  _Opcode_name[456:465],
	87:  _Opcode_name[465:473],
	88:  _Opcode_name[473:480],
	89:  _Opcode_name[480:487],
	90:  _Opcode_name[487:494],
	91:  _Opcode_name[494:501],
	92:  _Opcode_name[501:508],
	93:  _Opcode_name[508:515],
	94:  _Opcode_name[515:522],
	95:  _Opcode_name[522:528],
	96:  _Opcode_name[528:535],
	97:  _Opcode_name[535:542],
	98:  _Opcode_name[542:549],
	99:  _Opcode_name[549:556],
	100: _Opcode_name[556:563],
	101: _Opcode_name[563:570],
	102: _Opcode_name[570:577],
	103: _Opcode_name[577:583],
	104: _Opcode_name[583:589],
	105: _Opcode_name[589:595],
	106: _Opcode_name[595:601],
	107: _Opcode_name[601:607],
	108: _Opcode_name[607:613],
	109: _Opcode_name[613:619],
	110: _Opcode_name[619:625],
	111: _Opcode_name[625:630],
	112: _Opcode_name[630:636],
	113: _Opcode_name[636:642],
	114: _Opcode_name[642:648],
	115: _Opcode_name[648:654],
	116: _Opcode_name[654:660],
	117: _Opcode_name[660:666],
	118: _Opcode_name[666:672],
	119: _Opcode_name[672:677],
	120: _Opcode_name[677:683],
	121: _Opcode_name[683:689],
	122: _Opcode_name[689:695],
	123: _Opcode_name[695:701],
	124: _Opcode_name[701:707],
	125: _Opcode_name[707:713],
	126: _Opcode_name[713:719],
	127: _Opcode_name[719:724],
	128: _Opcode_name[724:730],
	129: _Opcode_name[730:736],
	130: _Opcode_name[736:742],
	131: _Opcode_name[742:748],
	132: _Opcode_name[748:754],
	133: _Opcode_name[754:760],
	134: _Opcode_name[760:766],
	135: _Opcode_name[766:771],
	136: _Opcode_name[771:780],
	137: _Opcode_name[780:786],
	139: _Opcode_name[786:789],
	140: _Opcode_name[789:795],
	141: _Opcode_name[795:799],
	142: _Opcode_name[799:804],
	144: _Opcode_name[804:810],
	145: _Opcode_name[810:813],
	146: _Opcode_name[813:815],
	147: _Opcode_name[815:818],
	151: _Opcode_name[818:823],
	152: _Opcode_name[823:831],
	153: _Opcode_name[831:835],
	154: _Opcode_name[835:838],
	155: _Opcode_name[838:844],
	156: _Opcode_name[844:847],
	157: _Opcode_name[847:850],
	158: _Opcode_name[850:853],
	159: _Opcode_name[853:856],
	160: _Opcode_name[856:859],
	161: _Opcode_name[859:862],
	162: _Opcode_name[862:865],
	163: _Opcode_name[865:868],
	164: _Opcode_name[868:872],
	168: _Opcode_name[872:875],
	169: _Opcode_name[875:878],
	170: _Opcode_name[878:881],
	171: _Opcode_name[881:888],
	172: _Opcode_name[888:894],
	177: _Opcode_name[894:896],
	179: _Opcode_name[896:904],
	180: _Opcode_name[904:915],
	181: _Opcode_name[915:917],
	182: _Opcode_name[917:919],
	183: _Opcode_name[919:921],
	184: _Opcode_name[921:923],
	185: _Opcode_name[923:926],
	186: _Opcode_name[926:929],
	187: _Opcode_name[929:935],
	192: _Opcode_name[935:939],
	193: _Opcode_name[939:945],
	194: _Opcode_name[945:954],
	195: _Opcode_name[954:962],
	196: _Opcode_name[962:972],
	197: _Opcode_name[972:982],
	198: _Opcode_name[982:991],
	200: _Opcode_name[991:997],
	202: _Opcode_name[997:1001],
	203: _Opcode_name[1001:1007],
	204: _Opcode_name[1007:1011],
	205: _Opcode_name[1011:1017],
	206: _Opcode_name[1017:1025],
	207: _Opcode_name[1025:1031],
	208: _Opcode_name[1031:1038],
	209: _Opcode_name[1038:1050],
	210: _Opcode_name[1050:1056],
	211: _Opcode_name[1056:1066],
	212: _Opcode_name[1066:1073],
	216: _Opcode_name[1073:1079],
	217: _Opcode_name[1079:1085],
	219: _Opcode_name[1085:1092],
}

func (i Opcode) String() string {
	if str, ok := _Opcode_map[i]; ok {
		return str
	}
	return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
}
