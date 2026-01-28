package entity

// Gender - جنسیت
type Gender uint8

const (
	GenderMale   Gender = 1 // مرد
	GenderFemale Gender = 2 // زن
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "مرد"
	case GenderFemale:
		return "زن"
	default:
		return "نامشخص"
	}
}

// RelationType - نوع رابطه (نسبت)
type RelationType uint8

const (
	RelationSelf    RelationType = 1 // خود بیمه‌شده
	RelationSpouse  RelationType = 2 // همسر
	RelationChild   RelationType = 3 // فرزند
	RelationParent  RelationType = 4 // والدین
	RelationSibling RelationType = 5 // خواهر/برادر
)

func (r RelationType) String() string {
	switch r {
	case RelationSelf:
		return "بیمه‌شده اصلی"
	case RelationSpouse:
		return "همسر"
	case RelationChild:
		return "فرزند"
	case RelationParent:
		return "والدین"
	case RelationSibling:
		return "خواهر/برادر"
	default:
		return "سایر"
	}
}

// ClaimType - نوع ادعا
type ClaimType uint8

const (
	ClaimTypeDrug                    ClaimType = 1  // داروخانه
	ClaimTypeHospitalization         ClaimType = 2  // بستری
	ClaimTypeDental                  ClaimType = 3  // دندانپزشکی
	ClaimTypeDoctorVisit             ClaimType = 4  // ویزیت
	ClaimTypeLabTest                 ClaimType = 5  // آزمایشگاه
	ClaimTypeImaging                 ClaimType = 6  // تصویربرداری
	ClaimTypePhysiotherapy           ClaimType = 7  // فیزیوتراپی
	ClaimTypeOutpatientSurgery       ClaimType = 8  // جراحی سرپایی
	ClaimTypeEmergency               ClaimType = 9  // اورژانس
	ClaimTypeMedicalEquipment        ClaimType = 10 // تجهیزات پزشکی
	ClaimTypeInjection               ClaimType = 13 // تزریقات
	ClaimTypeClinic                  ClaimType = 15 // سرپایی بیمارستان
)

func (c ClaimType) String() string {
	names := map[ClaimType]string{
		ClaimTypeDrug:              "داروخانه",
		ClaimTypeHospitalization:   "بستری",
		ClaimTypeDental:            "دندانپزشکی",
		ClaimTypeDoctorVisit:       "ویزیت",
		ClaimTypeLabTest:           "آزمایشگاه",
		ClaimTypeImaging:           "تصویربرداری",
		ClaimTypePhysiotherapy:     "فیزیوتراپی",
		ClaimTypeOutpatientSurgery: "جراحی سرپایی",
		ClaimTypeEmergency:         "اورژانس",
		ClaimTypeMedicalEquipment:  "تجهیزات پزشکی",
		ClaimTypeInjection:         "تزریقات",
		ClaimTypeClinic:            "سرپایی بیمارستان",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "نامشخص"
}

// ClaimStatus - وضعیت ادعا
type ClaimStatus uint8

const (
	ClaimStatusReturned          ClaimStatus = 1 // عودت شده
	ClaimStatusWaitRegister      ClaimStatus = 2 // منتظر ثبت
	ClaimStatusWaitCheck         ClaimStatus = 3 // منتظر ارزیابی
	ClaimStatusWaitCheckConfirm  ClaimStatus = 4 // منتظر تایید ارزیابی
	ClaimStatusWaitSendFinancial ClaimStatus = 5 // منتظر ارسال به مالی
	ClaimStatusArchived          ClaimStatus = 6 // پرداخت شده / آرشیو
	ClaimStatusWaitCheckAgain    ClaimStatus = 8 // ارزیابی مجدد
)

func (c ClaimStatus) String() string {
	names := map[ClaimStatus]string{
		ClaimStatusReturned:          "عودت شده",
		ClaimStatusWaitRegister:      "منتظر ثبت",
		ClaimStatusWaitCheck:         "منتظر ارزیابی",
		ClaimStatusWaitCheckConfirm:  "منتظر تایید",
		ClaimStatusWaitSendFinancial: "منتظر ارسال به مالی",
		ClaimStatusArchived:          "آرشیو",
		ClaimStatusWaitCheckAgain:    "ارزیابی مجدد",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "نامشخص"
}

// CenterType - نوع مرکز درمانی
type CenterType uint8

const (
	CenterTypeHospital      CenterType = 1 // بیمارستان
	CenterTypeClinic        CenterType = 2 // کلینیک/درمانگاه
	CenterTypePharmacy      CenterType = 3 // داروخانه
	CenterTypeLab           CenterType = 4 // آزمایشگاه
	CenterTypeImaging       CenterType = 5 // مرکز تصویربرداری
	CenterTypePhysiotherapy CenterType = 6 // فیزیوتراپی
	CenterTypeDental        CenterType = 7 // دندانپزشکی
	CenterTypeOptics        CenterType = 8 // عینک‌سازی
)

func (c CenterType) String() string {
	names := map[CenterType]string{
		CenterTypeHospital:      "بیمارستان",
		CenterTypeClinic:        "کلینیک/درمانگاه",
		CenterTypePharmacy:      "داروخانه",
		CenterTypeLab:           "آزمایشگاه",
		CenterTypeImaging:       "مرکز تصویربرداری",
		CenterTypePhysiotherapy: "فیزیوتراپی",
		CenterTypeDental:        "دندانپزشکی",
		CenterTypeOptics:        "عینک‌سازی",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "سایر"
}

// ContractStatus - وضعیت قرارداد
type ContractStatus uint8

const (
	ContractStatusActive   ContractStatus = 1 // فعال
	ContractStatusInactive ContractStatus = 2 // غیرفعال
	ContractStatusExpired  ContractStatus = 3 // منقضی
	ContractStatusPending  ContractStatus = 4 // در انتظار تایید
)

// PackageStatus - وضعیت بسته اسناد
type PackageStatus uint8

const (
	PackageStatusDraft            PackageStatus = 1 // پیش‌نویس
	PackageStatusWaitCheck        PackageStatus = 2 // منتظر ارزیابی
	PackageStatusWaitConfirm      PackageStatus = 3 // منتظر تایید
	PackageStatusWaitPayment      PackageStatus = 4 // منتظر پرداخت
	PackageStatusPaid             PackageStatus = 5 // پرداخت شده
	PackageStatusPartiallyPaid    PackageStatus = 6 // پرداخت ناقص
	PackageStatusReturned         PackageStatus = 7 // برگشت خورده
)

// AdmissionType - نوع پذیرش
type AdmissionType uint8

const (
	AdmissionOutpatient     AdmissionType = 1 // سرپایی
	AdmissionHospitalization AdmissionType = 2 // بستری
)

// DischargeCondition - وضعیت ترخیص
type DischargeCondition uint8

const (
	DischargeRecovered       DischargeCondition = 1 // بهبود یافته
	DischargeImproved        DischargeCondition = 2 // با حال بهتر
	DischargePersonalRequest DischargeCondition = 3 // درخواست شخصی
	DischargeTransferred     DischargeCondition = 4 // انتقالی
	DischargeDeceased        DischargeCondition = 5 // فوت شده
)
