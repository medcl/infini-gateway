package nlpt_detect

import "fmt"

func ExampleDetect() {
	lng := "नमोहन सिंह मंत्रिमंडल में मंत्रालयों का बँटवारा पूरा हो गया है."
	fmt.Println(Detect(lng, "name", len(lng), 3, 3, 3))
	// Output: HINDI
}

func ExampleStaticDetect() {
	fmt.Println(StaticDetect("Οἱ δὲ Φοίνιϰες οὗτοι οἱ σὺν Κάδμῳ ἀπιϰόμενοι")) // greek
	// Output: el
}

func ExampleGetLanguageName() {
	fmt.Println(GetLanguageName("筆記本在這邊｡")) //traditional chinese
	// Output: ChineseT
}

func ExampleGetLanguageDelcaredName() {
	fmt.Println(GetLanguageDeclaredName("人人生而自由，在尊严和权利上一律平等。他们赋有理性和良心，并应以兄弟关系的精神互相对待。")) //simplified chinese
	// Output: CHINESE
}

func ExampleGetLanguageCode() {
	fmt.Println(GetLanguageCode("Благодарим за письмо от")) // russian
	// Output: ru
}
