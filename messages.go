package main

import (
	"fmt"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
)

func messageYouAlreadySignedUp(name string) []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain(fmt.Sprintf("%s عزیز، شما قبلاً ثبت‌نام کرده‌اید.", name)),
	}
}

func messageAskFirstName() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً نام کوچک خود را به فارسی وارد کنید."),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("ابراهیم"),
	}
}

func messageAskLastName() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً نام خانوادگی خود را به فارسی وارد کنید."),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("نجاتی"),
	}
}

func messageAskPhone() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("لطفاً شماره تلفن خود را با اعداد لاتین وارد کنید."),
		styling.Plain("\n"),
		styling.Plain("برای مثال: "),
		styling.Bold("09123456789"),
	}
}

func messageStart() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("سلام. خوش اومدی!"),
	}
}

func messageIsNotPersian() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("فقط از حروف فارسی استفاده کنید."),
	}
}

func messageHasNoText() []message.StyledTextOption {
	return []message.StyledTextOption{
		styling.Plain("پیام شما حاوی متن نیست."),
	}
}
