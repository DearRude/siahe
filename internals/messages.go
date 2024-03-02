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

func MessageUserHasNoAccount() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("شما دارای حساب نیستید."),
	}
}
func MessageAccountDeleted() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("حساب شما حذف شد!"),
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

func MessageAddPlaceHelp() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("اطلاعات را با خط تیره - فاصله گذاری کنید. \n"),
		styling.Plain("اطلاعات مربوط به «مکان» عبارت‌اند از: "),
		styling.Bold("نام، آدرس، ظرفیت\n"),
		styling.Plain("برای مثال: \n\n"),
	}
}

func MessageAddPlaceExample() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("/addPlace"),
		styling.Plain("\n-\n"),
		styling.Plain("اتاق آروین"),
		styling.Plain("\n-\n"),
		styling.Plain("دانشگاه فردوسی، ساختمان کانون‌های فرهنگی، اتاق آروین"),
		styling.Plain("\n-\n"),
		styling.Plain("30"),
	}
}

func MessageAddEventeHelp() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("اطلاعات را با خط تیره - فاصله گذاری کنید. \n"),
		styling.Plain("اطلاعات مربوط به «رویداد» عبارت‌اند از: "),
		styling.Bold("نام، توضیحات، نیاز به هزینه؟، ماکسیمم خرید یکجا بلیط، آیدی مکان رویداد\n"),
		styling.Plain("برای مثال: \n\n"),
	}
}

func MessageAddEventExample() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("/addEvent"),
		styling.Plain("\n-\n"),
		styling.Plain("نشست نمایش‌نامه خوانی چهارصندوق"),
		styling.Plain("\n-\n"),
		styling.Plain("در این نشست به متن‌خوانی نمایش‌نامه چهارصندوق می‌پردازیم. شرکت در این نشست برای دانشجویان رایگان و آزاد است."),
		styling.Plain("\n-\n"),
		styling.Plain("بله"),
		styling.Plain("\n-\n"),
		styling.Plain("5"),
		styling.Plain("\n-\n"),
		styling.Plain("2"),
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

	m := []message.StyledTextOption{
		styling.Bold("کاربر با آیدی: "),
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

func MessageIsUserInfoCorrect(user db.User) []message.StyledTextOption {
	m := []message.StyledTextOption{
		styling.Plain("آیا اطلاعات وارد شده صحیح است؟"),
		styling.Plain("\n\n"), // newline
	}

	return append(m, MessagePrintUser(user)...)
}

func MessagePrintPlace(place db.Place) []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Bold("مکان با آیدی: "),
		styling.Code(fmt.Sprintf("%d", place.ID)),
		styling.Plain("\n"), // newline
		styling.Bold("نام: "),
		styling.Plain(place.Name),
		styling.Plain("\n"), // newline
		styling.Bold("آدرس: "),
		styling.Plain(place.Address),
		styling.Plain("\n"), // newline
		styling.Bold("ظرفیت: "),
		styling.Code(fmt.Sprintf("%d", place.Capacity)),
	}
}

func MessagePrintPlaces(places []db.Place) []message.StyledTextOption {
	m := []message.StyledTextOption{
		styling.Bold("مکان‌های ثبت شده: \n\n"),
	}

	for _, place := range places {
		m = append(m, []styling.StyledTextOption{
			styling.Code(fmt.Sprintf("%d", place.ID)),
			styling.Plain(fmt.Sprintf(": %s\n", place.Name)),
		}...)
	}

	return m
}

func MessagePlaceAdded(place db.Place) []message.StyledTextOption {
	m := []message.StyledTextOption{
		styling.Plain("مکان زیر اضافه شد!"),
		styling.Plain("\n\n"), // newline
	}

	return append(m, MessagePrintPlace(place)...)
}

func MessagePrintEvent(event db.Event) []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Bold("رویداد با آیدی: "),
		styling.Code(fmt.Sprintf("%d", event.ID)),
		styling.Plain("\n"), // newline
		styling.Bold("نام: "),
		styling.Plain(event.Name),
		styling.Plain("\n"), // newline
		styling.Bold("توضیحات: "),
		styling.Plain(event.Description),
		styling.Plain("\n"), // newline
		styling.Bold("نیاز به هزینه؟ "),
		styling.Plain(boolToText(event.IsPaid)),
		styling.Plain("\n"), // newline
		styling.Bold("ماکسیمم خرید یک‌جا بلیط: "),
		styling.Code(fmt.Sprintf("%d", event.MaxTicketBatch)),
		styling.Plain("\n"), // newline
		styling.Bold("آیدی مکان رویداد: "),
		styling.Code(fmt.Sprintf("%d", event.PlaceID)),
		styling.Plain("\n"), // newline
		styling.Bold("فعال است؟ "),
		styling.Plain(boolToText(event.IsActive)),
	}
}

func MessagePrintEvents(events []db.Event) []message.StyledTextOption {
	m := []message.StyledTextOption{
		styling.Bold("رویدادهای ثبت شده: \n\n"),
	}

	for _, event := range events {
		m = append(m, []styling.StyledTextOption{
			styling.Code(fmt.Sprintf("%d", event.ID)),
			styling.Plain(fmt.Sprintf(": %s\n", event.Name)),
		}...)
	}

	return m
}

func MessageEventAdded(event db.Event) []message.StyledTextOption {
	m := []message.StyledTextOption{
		styling.Plain("رویداد زیر اضافه شد!"),
		styling.Plain("\n\n"), // newline
	}

	return append(m, MessagePrintEvent(event)...)
}

func boolToText(b bool) string {
	if b {
		return "بله"
	}
	return "خیر"
}

func boolToGender(b bool) string {
	if b {
		return "مرد"
	}
	return "زن"
}
