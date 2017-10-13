package x86

var vexOptab = []Optab{
	{AANDNL, yvex_r3, Pvex, [23]uint8{VEX_NDS_LZ_0F38_W0, 0xF2}},
	{AANDNQ, yvex_r3, Pvex, [23]uint8{VEX_NDS_LZ_0F38_W1, 0xF2}},
	{ABEXTRL, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_0F38_W0, 0xF7}},
	{ABEXTRQ, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_0F38_W1, 0xF7}},
	{ABZHIL, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_0F38_W0, 0xF5}},
	{ABZHIQ, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_0F38_W1, 0xF5}},
	{AMULXL, yvex_r3, Pvex, [23]uint8{VEX_NDD_LZ_F2_0F38_W0, 0xF6}},
	{AMULXQ, yvex_r3, Pvex, [23]uint8{VEX_NDD_LZ_F2_0F38_W1, 0xF6}},
	{APDEPL, yvex_r3, Pvex, [23]uint8{VEX_NDS_LZ_F2_0F38_W0, 0xF5}},
	{APDEPQ, yvex_r3, Pvex, [23]uint8{VEX_NDS_LZ_F2_0F38_W1, 0xF5}},
	{APEXTL, yvex_r3, Pvex, [23]uint8{VEX_NDS_LZ_F3_0F38_W0, 0xF5}},
	{APEXTQ, yvex_r3, Pvex, [23]uint8{VEX_NDS_LZ_F3_0F38_W1, 0xF5}},
	{ASARXL, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_F3_0F38_W0, 0xF7}},
	{ASARXQ, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_F3_0F38_W1, 0xF7}},
	{ASHLXL, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_66_0F38_W0, 0xF7}},
	{ASHLXQ, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_66_0F38_W1, 0xF7}},
	{ASHRXL, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_F2_0F38_W0, 0xF7}},
	{ASHRXQ, yvex_vmr3, Pvex, [23]uint8{VEX_NDS_LZ_F2_0F38_W1, 0xF7}},
	{AVZEROUPPER, ynone, Px, [23]uint8{0xc5, 0xf8, 0x77}},
	{AVMOVDQU, yvex_vmovdqa, Pvex, [23]uint8{VEX_NOVSR_128_F3_0F_WIG, 0x6F, VEX_NOVSR_128_F3_0F_WIG, 0x7F, VEX_NOVSR_256_F3_0F_WIG, 0x6F, VEX_NOVSR_256_F3_0F_WIG, 0x7F}},
	{AVMOVDQA, yvex_vmovdqa, Pvex, [23]uint8{VEX_NOVSR_128_66_0F_WIG, 0x6F, VEX_NOVSR_128_66_0F_WIG, 0x7F, VEX_NOVSR_256_66_0F_WIG, 0x6F, VEX_NOVSR_256_66_0F_WIG, 0x7F}},
	{AVMOVNTDQ, yvex_vmovntdq, Pvex, [23]uint8{VEX_NOVSR_128_66_0F_WIG, 0xE7, VEX_NOVSR_256_66_0F_WIG, 0xE7}},
	{AVPCMPEQB, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0x74, VEX_NDS_256_66_0F_WIG, 0x74}},
	{AVPXOR, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0xEF, VEX_NDS_256_66_0F_WIG, 0xEF}},
	{AVPMOVMSKB, yvex_xyr2, Pvex, [23]uint8{VEX_NOVSR_128_66_0F_WIG, 0xD7, VEX_NOVSR_256_66_0F_WIG, 0xD7}},
	{AVPAND, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0xDB, VEX_NDS_256_66_0F_WIG, 0xDB}},
	{AVPBROADCASTB, yvex_vpbroadcast, Pvex, [23]uint8{VEX_NOVSR_128_66_0F38_W0, 0x78, VEX_NOVSR_256_66_0F38_W0, 0x78}},
	{AVPTEST, yvex_xy2, Pvex, [23]uint8{VEX_NOVSR_128_66_0F38_WIG, 0x17, VEX_NOVSR_256_66_0F38_WIG, 0x17}},
	{AVPSHUFB, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F38_WIG, 0x00, VEX_NDS_256_66_0F38_WIG, 0x00}},
	{AVPSHUFD, yvex_xyi3, Pvex, [23]uint8{VEX_NOVSR_128_66_0F_WIG, 0x70, VEX_NOVSR_256_66_0F_WIG, 0x70, VEX_NOVSR_128_66_0F_WIG, 0x70, VEX_NOVSR_256_66_0F_WIG, 0x70}},
	{AVPOR, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0xeb, VEX_NDS_256_66_0F_WIG, 0xeb}},
	{AVPADDQ, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0xd4, VEX_NDS_256_66_0F_WIG, 0xd4}},
	{AVPADDD, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0xfe, VEX_NDS_256_66_0F_WIG, 0xfe}},
	{AVADDSD, yvex_x3, Pvex, [23]uint8{VEX_NDS_128_F2_0F_WIG, 0x58}},
	{AVSUBSD, yvex_x3, Pvex, [23]uint8{VEX_NDS_128_F2_0F_WIG, 0x5c}},
	{AVFMADD213SD, yvex_x3, Pvex, [23]uint8{VEX_DDS_LIG_66_0F38_W1, 0xa9}},
	{AVFMADD231SD, yvex_x3, Pvex, [23]uint8{VEX_DDS_LIG_66_0F38_W1, 0xb9}},
	{AVFNMADD213SD, yvex_x3, Pvex, [23]uint8{VEX_DDS_LIG_66_0F38_W1, 0xad}},
	{AVFNMADD231SD, yvex_x3, Pvex, [23]uint8{VEX_DDS_LIG_66_0F38_W1, 0xbd}},
	{AVPSLLD, yvex_shift, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0x72, 0xf0, VEX_NDS_256_66_0F_WIG, 0x72, 0xf0, VEX_NDD_128_66_0F_WIG, 0xf2, VEX_NDD_256_66_0F_WIG, 0xf2}},
	{AVPSLLQ, yvex_shift, Pvex, [23]uint8{VEX_NDD_128_66_0F_WIG, 0x73, 0xf0, VEX_NDD_256_66_0F_WIG, 0x73, 0xf0, VEX_NDS_128_66_0F_WIG, 0xf3, VEX_NDS_256_66_0F_WIG, 0xf3}},
	{AVPSRLD, yvex_shift, Pvex, [23]uint8{VEX_NDD_128_66_0F_WIG, 0x72, 0xd0, VEX_NDD_256_66_0F_WIG, 0x72, 0xd0, VEX_NDD_128_66_0F_WIG, 0xd2, VEX_NDD_256_66_0F_WIG, 0xd2}},
	{AVPSRLQ, yvex_shift, Pvex, [23]uint8{VEX_NDD_128_66_0F_WIG, 0x73, 0xd0, VEX_NDD_256_66_0F_WIG, 0x73, 0xd0, VEX_NDS_128_66_0F_WIG, 0xd3, VEX_NDS_256_66_0F_WIG, 0xd3}},
	{AVPSRLDQ, yvex_shift_dq, Pvex, [23]uint8{VEX_NDD_128_66_0F_WIG, 0x73, 0xd8, VEX_NDD_256_66_0F_WIG, 0x73, 0xd8}},
	{AVPSLLDQ, yvex_shift_dq, Pvex, [23]uint8{VEX_NDD_128_66_0F_WIG, 0x73, 0xf8, VEX_NDD_256_66_0F_WIG, 0x73, 0xf8}},
	{AVPERM2F128, yvex_yyi4, Pvex, [23]uint8{VEX_NDS_256_66_0F3A_W0, 0x06}},
	{AVPALIGNR, yvex_yyi4, Pvex, [23]uint8{VEX_NDS_256_66_0F3A_WIG, 0x0f}},
	{AVPBLENDD, yvex_yyi4, Pvex, [23]uint8{VEX_NDS_256_66_0F3A_WIG, 0x02}},
	{AVINSERTI128, yvex_xyi4, Pvex, [23]uint8{VEX_NDS_256_66_0F3A_WIG, 0x38}},
	{AVPERM2I128, yvex_yyi4, Pvex, [23]uint8{VEX_NDS_256_66_0F3A_WIG, 0x46}},
	{ARORXL, yvex_ri3, Pvex, [23]uint8{VEX_NOVSR_LZ_F2_0F3A_W0, 0xf0}},
	{ARORXQ, yvex_ri3, Pvex, [23]uint8{VEX_NOVSR_LZ_F2_0F3A_W1, 0xf0}},
	{AVBROADCASTSD, yvex_vpbroadcast_sd, Pvex, [23]uint8{VEX_NOVSR_256_66_0F38_W0, 0x19}},
	{AVBROADCASTSS, yvex_vpbroadcast, Pvex, [23]uint8{VEX_NOVSR_128_66_0F38_W0, 0x18, VEX_NOVSR_256_66_0F38_W0, 0x18}},
	{AVMOVDDUP, yvex_xy2, Pvex, [23]uint8{VEX_NOVSR_128_F2_0F_WIG, 0x12, VEX_NOVSR_256_F2_0F_WIG, 0x12}},
	{AVMOVSHDUP, yvex_xy2, Pvex, [23]uint8{VEX_NOVSR_128_F3_0F_WIG, 0x16, VEX_NOVSR_256_F3_0F_WIG, 0x16}},
	{AVMOVSLDUP, yvex_xy2, Pvex, [23]uint8{VEX_NOVSR_128_F3_0F_WIG, 0x12, VEX_NOVSR_256_F3_0F_WIG, 0x12}},
}
