package internals

import (
	"fmt"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"

	db "github.com/DearRude/fumTheatreBot/database"
)

func MessageYouAlreadySignedUp(name string) []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain(fmt.Sprintf("%s عزیز، شما قبلاً ثبت‌نام کرده‌اید.", name)),
	}
}

func MessageAskFirstName() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً نام کوچک خود را به فارسی وارد کنید."),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("ابراهیم"),
	}
}

func MessageAskLastName() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً نام خانوادگی خود را به فارسی وارد کنید."),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("نجاتی"),
	}
}

func MessageAskPhone() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً شماره تلفن خود را با زدن دکمه Markup وارد کنید."),
	}
}

func MessageStart(id int64) []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("سلام. خوش اومدی!\n"),
		styling.Plain("این آیدی تلگرام توئه: "),
		styling.Code(fmt.Sprintf("%d", id)),
	}
}

func MessageIsNotPersian() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("فقط از حروف فارسی استفاده کنید."),
	}
}

func MessageIsNotEntranceYear() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً سال تحصیلی خود را به صورت صحیح وارد کنید."),
	}
}

func MessageIsNotPhone() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً شماره تلفن خود را به صورت صحیح وارد کنید."),
	}
}

func MessageIsNotStudentNumber() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً شماره دانشجویی خود را به صورت صحیح وارد کنید."),
	}
}

func MessageHasNoText() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("پیام شما حاوی متن نیست."),
	}
}

func MessageAskGender() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("جنسیت خود را انتخاب کنید:"),
	}
}

func MessageAskIsFUMStudent() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("آیا هم‌اکنون دانشجوی دانشگاه فردوسی مشهد هستید؟"),
	}
}

func MessageAskStudentNumber() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("شماره دانشجویی خود را با اعداد لاتین وارد کنید:"),
	}
}

func MessageAskFumFaculty() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("دانشکده تحصیل خود را انتخاب کنید:"),
	}
}

func MessageAskIsStudent() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("آیا اکنون دانشجو هستید؟"),
	}
}

func MessageAskIsMashhad() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("آیا دانشجوی یکی از دانشگاه‌های مشهد هستید؟"),
	}
}

func MessageAskIsMastPhd() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("آیا دانشجوی تحصیلات تکمیلی هستید؟"),
	}
}

func MessageAskUniversityName() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("نام دانشگاه خود را به فارسی وارد کنید:"),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("دانشگاه صنعتی اصفهان"),
	}
}

func MessageAskMajor() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("رشته تحصیلی خود را به فارسی وارد کنید:"),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("مهندسی عمران"),
	}
}

func MessageAskEntranceYear() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("سال ورودی به دانشگاه را با اعداد لاتین وارد کنید:"),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("1398"),
	}
}

func MessageAskIsGraduate() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("آیا دانشجوی فارغ تحصیل دانشگاه فردوسی مشهد هستید؟"),
	}
}

func MessageAskIsRelative() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("آیا خانواده درجه اول دانشجویان یا از اساتید و کارکنان دانشگاه فردوسی مشهد هستید؟"),
	}
}

func MessageSignUpFinished(name string) []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain(fmt.Sprintf("%s جان، ثبت نام شما به پایان رسید.", name)),
	}
}

func MessageCancelSignUp() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("ثبت‌نام کنسل شد. در صورت درخواست، دوباره امتحان کنید"),
	}
}

func ButtonYesNo() []tg.KeyboardButtonClass {
	return []tg.KeyboardButtonClass{
		markup.Callback("بله", []byte("yes")),
		markup.Callback("خیر", []byte("no")),
	}
}

func ButtonAskGender() []tg.KeyboardButtonClass {
	return []tg.KeyboardButtonClass{
		&tg.KeyboardButtonCallback{
			Text: "آقا",
			Data: []byte("boy"),
		},
		&tg.KeyboardButtonCallback{
			Text: "خانم",
			Data: []byte("girl"),
		},
	}
}

func ButtonAskFumFaculty() tg.ReplyMarkupClass {
	callbackSameName := func(text string) *tg.KeyboardButtonCallback {
		return markup.Callback(text, []byte(text))
	}

	return markup.InlineKeyboard(
		markup.Row(
			callbackSameName("ادبیات و علوم انسانی"),
			callbackSameName("الهیات و معارف اسلامی"),
		),
		markup.Row(
			callbackSameName("حقوق و علوم سیاسی"),
			callbackSameName("دامپزشکی"),
		),
		markup.Row(
			callbackSameName("علوم"),
			callbackSameName("علوم اداری و اقتصادی"),
		),
		markup.Row(
			callbackSameName("علوم تربیتی و روانشناسی"),
			callbackSameName("علوم ریاضی"),
		),
		markup.Row(
			callbackSameName("علوم ورزشی"),
			callbackSameName("کشاورزی"),
		),
		markup.Row(
			callbackSameName("معماری و شهرسازی"),
			callbackSameName("منابع طبیعی و محیط زیست"),
		),
		markup.Row(
			callbackSameName("مهندسی"),
			callbackSameName("هنر نیشابور"),
		),
	)
}

func ButtonAskPhone() tg.ReplyMarkupClass {
	return markup.BuildKeyboard().
		Resize().
		SingleUse().
		Build(markup.Row(markup.
			RequestPhone("شماره خود را به اشتراک بگذارید")))
}

func MessagePrintUser(user db.User) []message.StyledTextOption {
	boolToText := func(b bool) string {
		if b {
			return "بله"
		}
		return "خیر"
	}

	boolToGender := func(b bool) string {
		if b {
			return "مرد"
		}
		return "زن"
	}

	m := []message.StyledTextOption{
		styling.Plain("کاربر با آیدی: "),
		styling.Code(fmt.Sprintf("%d", user.ID)),
		styling.Plain("\n"), // newline
		styling.Bold("نام: "),
		styling.Plain(user.FirstName),
		styling.Plain("\n"), // newline
		styling.Bold("نام خانوادگی: "),
		styling.Plain(user.LastName),
		styling.Plain("\n"), // newline
		styling.Bold("شماره تلفن: "),
		styling.Phone(user.PhoneNumber),
		styling.Plain("\n"), // newline
		styling.Bold("جنس: "),
		styling.Plain(boolToGender(user.IsBoy)),
		styling.Plain("\n"), // newline
		styling.Bold("دانشجوی فردوسی؟ "),
		styling.Plain(boolToText(user.IsFumStudent)),
		styling.Plain("\n"), // newline
		styling.Bold("دانشجو؟ "),
		styling.Plain(boolToText(user.IsStudent)),
		styling.Plain("\n"), // newline
		styling.Bold("دانشجوی مشهد؟ "),
		styling.Plain(boolToText(user.IsMashhadStudent)),
		styling.Plain("\n"), // newline
		styling.Bold("دانشجوی فارغ التحصیل؟ "),
		styling.Plain(boolToText(user.IsGraduateStudent)),
		styling.Plain("\n"), // newline
		styling.Bold("خانواده فردوسی؟ "),
		styling.Plain(boolToText(user.IsStudentRelative)),
		styling.Plain("\n"), // newline
	}
	if user.IsStudent {
		m = append(m, []styling.StyledTextOption{
			styling.Bold("نام دانشگاه: "),
			styling.Plain(user.UniversityName),
			styling.Plain("\n"), // newline
			styling.Bold("سال ورود: "),
			styling.Code(user.EntranceYear),
			styling.Plain("\n"), // newline
			styling.Bold("دانشجوی تحصیلات تکمیلی؟ "),
			styling.Plain(boolToText(user.IsMashhadStudent)),
			styling.Plain("\n"), // newline
			styling.Bold("رشته تحصیلی: "),
			styling.Plain(user.StudentMajor),
			styling.Plain("\n"), // newline
		}...)
	}
	if user.IsFumStudent {
		m = append(m, []styling.StyledTextOption{
			styling.Bold("شماره دانشجویی: "),
			styling.Code(user.StudentNumber),
			styling.Plain("\n"), // newline
			styling.Bold("دانشکده تحصیل: "),
			styling.Plain(user.FumFaculty),
		}...)
	}

	return m
}

func MessageIsInfoCorrect(user db.User) []message.StyledTextOption {

	m := []message.StyledTextOption{
		styling.Plain("آیا اطلاعات وارد شده صحیح است؟"),
		styling.Plain("\n\n"), // newline
	}

	return append(m, MessagePrintUser(user)...)
}
