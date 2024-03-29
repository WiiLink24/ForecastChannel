package main

var countryCodes = map[string]uint8{}

func PopulateCountryCodes() {
	countryCodes["Japan"] = 1
	countryCodes["Antarctica"] = 2
	countryCodes["Caribbean Netherlands"] = 3
	countryCodes["Falkland Islands"] = 4
	countryCodes["Sint Maarten"] = 7
	countryCodes["Australia"] = 65
	countryCodes["Anguilla"] = 8
	countryCodes["Antigua and Barbuda"] = 9
	countryCodes["Argentina"] = 10
	countryCodes["Aruba"] = 11
	countryCodes["Bahamas"] = 12
	countryCodes["Barbados"] = 13
	countryCodes["Belize"] = 14
	countryCodes["Bolivia"] = 15
	countryCodes["Brazil"] = 16
	countryCodes["British Virgin Islands"] = 17
	countryCodes["Canada"] = 18
	countryCodes["Cayman Islands"] = 19
	countryCodes["Chile"] = 20
	countryCodes["Colombia"] = 21
	countryCodes["Costa Rica"] = 22
	countryCodes["Curaçao"] = 38
	countryCodes["Dominica"] = 23
	countryCodes["Dominican Republic"] = 24
	countryCodes["Ecuador"] = 25
	countryCodes["El Salvador"] = 26
	countryCodes["French Guiana"] = 27
	countryCodes["Grenada"] = 28
	countryCodes["Guadeloupe"] = 29
	countryCodes["Guatemala"] = 30
	countryCodes["Guyana"] = 31
	countryCodes["Haiti"] = 32
	countryCodes["Honduras"] = 33
	countryCodes["Jamaica"] = 34
	countryCodes["Martinique"] = 35
	countryCodes["Mexico"] = 36
	countryCodes["Montserrat"] = 37
	countryCodes["Nicaragua"] = 39
	countryCodes["Panama"] = 40
	countryCodes["Paraguay"] = 41
	countryCodes["Peru"] = 42
	countryCodes["St. Kitts and Nevis"] = 43
	countryCodes["St. Lucia"] = 44
	countryCodes["St. Vincent and the Grenadines"] = 45
	countryCodes["Suriname"] = 46
	countryCodes["Trinidad and Tobago"] = 47
	countryCodes["Turks and Caicos Islands"] = 48
	countryCodes["United States"] = 49
	countryCodes["Uruguay"] = 50
	countryCodes["US Virgin Islands"] = 51
	countryCodes["Venezuela"] = 52
	countryCodes["Armenia"] = 53
	countryCodes["Belarus"] = 54
	countryCodes["Georgia"] = 55
	countryCodes["Kosovo"] = 56
	countryCodes["Faroe Islands"] = 63
	countryCodes["Albania"] = 64
	countryCodes["Australia"] = 65
	countryCodes["Austria"] = 66
	countryCodes["Belgium"] = 67
	countryCodes["Bosnia and Herzegovina"] = 68
	countryCodes["Botswana"] = 69
	countryCodes["Bulgaria"] = 70
	countryCodes["Croatia"] = 71
	countryCodes["Cyprus"] = 72
	countryCodes["Czechia"] = 73
	countryCodes["Denmark"] = 74
	countryCodes["Estonia"] = 75
	countryCodes["Finland"] = 76
	countryCodes["France"] = 77
	countryCodes["Germany"] = 78
	countryCodes["Greece"] = 79
	countryCodes["Hungary"] = 80
	countryCodes["Iceland"] = 81
	countryCodes["Ireland"] = 82
	countryCodes["Italy"] = 83
	countryCodes["Latvia"] = 84
	countryCodes["Lesotho"] = 85
	countryCodes["Liechtenstein"] = 86
	countryCodes["Lithuania"] = 87
	countryCodes["Luxembourg"] = 88
	countryCodes["North Macedonia"] = 89
	countryCodes["Malta"] = 90
	countryCodes["Montenegro"] = 91
	countryCodes["Mozambique"] = 92
	countryCodes["Namibia"] = 93
	countryCodes["Netherlands"] = 94
	countryCodes["New Zealand"] = 95
	countryCodes["Norway"] = 96
	countryCodes["Poland"] = 97
	countryCodes["Portugal"] = 98
	countryCodes["Romania"] = 99
	countryCodes["Russia"] = 100
	countryCodes["Serbia"] = 101
	countryCodes["Slovakia"] = 102
	countryCodes["Slovenia"] = 103
	countryCodes["South Africa"] = 104
	countryCodes["Spain"] = 105
	countryCodes["Eswatini"] = 106
	countryCodes["Sweden"] = 107
	countryCodes["Switzerland"] = 108
	countryCodes["Türkiye"] = 109
	countryCodes["United Kingdom"] = 110
	countryCodes["Zambia"] = 111
	countryCodes["Zimbabwe"] = 112
	countryCodes["Azerbaijan"] = 113
	countryCodes["Mauritania"] = 114
	countryCodes["Mali"] = 115
	countryCodes["Niger"] = 116
	countryCodes["Chad"] = 117
	countryCodes["Sudan"] = 118
	countryCodes["Eritrea"] = 119
	countryCodes["Dijibouti"] = 120
	countryCodes["Somalia"] = 121
	countryCodes["Andorra"] = 122
	countryCodes["Guernsey"] = 124
	countryCodes["Isle of Man"] = 125
	countryCodes["Jersey"] = 126
	countryCodes["Monaco"] = 127
	countryCodes["Taiwan"] = 128
	countryCodes["Cambodia"] = 129
	countryCodes["Laos"] = 130
	countryCodes["Mongolia"] = 131
	countryCodes["Myanmar"] = 132
	countryCodes["Nepal"] = 133
	countryCodes["Vietnam"] = 134
	countryCodes["North Korea"] = 135
	countryCodes["South Korea"] = 136
	countryCodes["Bangladesh"] = 137
	countryCodes["Bhutan"] = 138
	countryCodes["Brunei"] = 139
	countryCodes["Maldives"] = 140
	countryCodes["Sri Lanka"] = 141
	countryCodes["East Timor"] = 142
	countryCodes["British Indian Ocean Territory"] = 143
	countryCodes["Hong Kong"] = 144
	countryCodes["Macao"] = 145
	countryCodes["Cook Islands"] = 146
	countryCodes["Niue"] = 147
	countryCodes["Northern Mariana Islands"] = 149
	countryCodes["American Samoa"] = 150
	countryCodes["Guam"] = 151
	countryCodes["Indonesia"] = 152
	countryCodes["Singapore"] = 153
	countryCodes["Thailand"] = 154
	countryCodes["Philippines"] = 155
	countryCodes["Malaysia"] = 156
	countryCodes["Saint Barthélemy"] = 157
	countryCodes["Saint Martin"] = 158
	countryCodes["Saint Pierre and Miquelon"] = 159
	countryCodes["China"] = 160
	countryCodes["Afghanistan"] = 161
	countryCodes["Kazakhstan"] = 162
	countryCodes["Kyrgyzstan"] = 163
	countryCodes["Pakistan"] = 164
	countryCodes["Tajikistan"] = 165
	countryCodes["Turkmenistan"] = 166
	countryCodes["Uzbekistan"] = 167
	countryCodes["United Arab Emirates"] = 168
	countryCodes["India"] = 169
	countryCodes["Egypt"] = 170
	countryCodes["Oman"] = 171
	countryCodes["Qatar"] = 172
	countryCodes["Kuwait"] = 173
	countryCodes["Saudi Arabia"] = 174
	countryCodes["Syria"] = 175
	countryCodes["Bahrain"] = 176
	countryCodes["Jordan"] = 177
	countryCodes["Iran"] = 178
	countryCodes["Iraq"] = 179
	countryCodes["Israel"] = 180
	countryCodes["Moldova"] = 207
	countryCodes["Ukraine"] = 208
	countryCodes["Libya"] = 218
	countryCodes["Morocco"] = 219
	countryCodes["South Sudan"] = 220
	countryCodes["Cuba"] = 223
	countryCodes["Kiribati"] = 240
}
