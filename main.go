package main

import (
	"fmt"
	"translator/openai"
)

func main() {
	// fmt.Println("Hello, world.")
	// textMsg := "حلو، أنا يعني حابب بس أأكد على علا قاله في حِتة الـ BDM، فده الراجل اللي بيجيب الـ opportunities أو عمال يجيب leads وبعد كده بيروح للـ sales أو الـ opportunity owner، أو يعني تقدر تقول هو اللي بيمانج الـ opportunity. بس الراجل ده في الشركات الـ corporate بيبقى عنده quota محملة بمنتجات كتير جدًا، تمام، وبيبقى عنده مشاريع كتير جدًا وبقى عنده detailed design، ودي محتاج يكفرها ومحتاج يغطيها. وبيبقى في questions كتير بيسألها، وبيبقى في details محتاجة تتشرح وبيبقى في محتاجة تتعمل لأن أنا في الآخر لما ببقى اشتغل في شركة tech، فالـ stakeholders بتاعي، أه طبعًا في managers وفي executives بس معظم الناس اللي أنا محتاج أقنعهم هما يا إما developers، يا إما technical managers، يا إما engineers نفسهم. ففيه faces عندي في الأوضة، فده لازم الناس دي لازم أكفر، لازم يعني أرضي الناس دي. فإزاي أرضيهم؟ إزاي إن أنا أقدر أديلهم المحتوى اللي هم عايزين يقتنعوا بيه؟ إن إزاي الحل أو المنتج بتاعي يناسب الشركة بتاعتهم، أو يناسب المشكلة اللي بتواجههم، أو يناسب وجهة النظر اللي هم عايزين يتجهوا ليها كمان سنتين أو أربع سنين. فهو ده الـ technical sales. فمعظم شغلنا، لو هقول مكونات الشغل، معظم شغلي إما بعمل حاجة اسمها deep dive presentation، إن أنا بخش، السيلز بيبقى عنده بريزنتيشن عامة، أنا ببدأ أخش في details، bucket بتاعته وfeatures بتاعته، وبعد كده أعمل حاجة اسمها demo، إن أنا أبدأ أعرض الـ product بتاعي وأبدأ أشوف، أوري الفيتشر بتاعته، هنا دي هيلها بعدين، وبعد كده بقى في طبعًا Q&A، ويعني حاجات كتير جدًا السيلز بتغطيها. فده بس يعني باختصار شديد الفرق ما بين السيلز اللي هو بيتكلم business وبيتكلم value، إنما الـ pre-sales بيتكلم technical. في الآخر هما الاتنين في الآخر الناس كلها قاعدة في نفس الأوضة بس كل واحد بيتكلم مع ناس مختلفة. الناس محتاجة تبقى فاهمه هو أصلاً الـ sales لو هو من الجزء بتاع الـ tech بيعمل إيه؟ يعني بيبيع إزاي؟ الـ process يعني، أوكي. طيب هو الـ sales في شركات الـ tech هم الناس اللي بيسموهم quota carriers، يعني الناس اللي عندهم التارجت بتاع البيع للـ solution اللي بتعمله الشركة دي سواء الـ solution ده كان software، hardware، أو حتى ممكن يكون services. ممكن تكون شركة software house يعني عارف شغلتهم الحقيقية إن هم بيعملوا code وبيطوروا applications بناءً على طلب الزباين بتوعهم. شغلة البياع إن هو يدور على الفرص اللي موجودة في الـ market وياخد الفرص دي يحولها إلى opportunities حقيقية ويمشي معاها من أول ما لقاها لحد ما يوصل customer في الآخر إن هو يطلع الـ PO أو يمضي على العقد ويؤمّت على الـ business. الشركات الشاطرة، الصغيرة لحد الكبيرة، هي الشركات اللي البياعين فيها عندهم authority كويسة إن هو يقدر يدور على الـ opportunities بشكل مضبوط ويقود الـ process ده كله من أوله لآخره، ويكون عنده target-driven. يعني بنشوف ساعات في شركات صغيرة بالذات startups يقولك: إحنا عندنا ناس business developers. الـ business developer ده مش هيخليك تقدر تطلع من حتة لحتة، لأن الـ sales هي شغلانة driven باي بشكل رئيسي الأرقام. فأنت لو بتعين شخص الأرقام مش بتعني له حاجة، مش بيحس إن هو عنده تحدي في إنه يحقق التارجت بتاعه، الشخص ده غالبًا مش هيقدر يبقى sales كويس، ممكن يشتغل business developer., انا هفاجئك شويه في الnegotiation skills دي. بس عايز اقوللك على نقطتين تانيين كويسين مهمين برضو شغلات sales وpre-sales  عارف انت في جمله كده كليشيه بتتكتب في الjob description بتاع اي jobs موجوده على sales وpre-sales يقوللك ده self-starter حلو جدا ده الجمله دي على قد ما هي كليشيه جدا بس هي يعني بتعكس الحقيقه جدا انت تصحى من نوم كل يوم وانت بتشتغل في الsales او الpre-sales لو ما بدات الحاجه هي مش تبدا لو ما خلصتش الحاجه هي مش هتخلص مافيش team عارف انت بتشتغل انت والشخص الpre-sales ده انت انت لوحدك يعني عارف انت انت لو ما بدات مافيش حاجه لو ما خلصتها مافيش حاجه بتخلص وهي دي المشكله هي شويه الlonely job يعني لازم تبقى متعود انك تقوم وتحرك نفسك لوحدك مش هتلاقي team بيعمل لك task وال والا انت تلاقي نفسك قاعد مافيش حاجه تعملها  يعني عارف هو competition كتير والشركات بtechnology كتير وكلهم بيجروا ورا نفس مجموعه الcustomers انت في الاخر وانت بتشتغل sales بتبقى focused على حته معينه سواء حته جغرافيه او vertical مثلا banking public sector مش عارف ايه فانت كل الشركات بتجري ورا نفس الناس هو الراجل اللي انت عايز تكلمه ده في 50 غيرك عايزين يكلمه فانت لو ما بدات ال negotiation ما بتبدأ. جابر قاللك من شويه جزء مهم جدا من الشغلانه انك تعرف تسال الاسئله الصح ما تسالش الاسئلها للي مالهاش مالهاش قيمه اللي هي مش هتضيفلك حاجه فانت اهم اهم سؤال هو الشخص اللي عايز يشتري منك ده عايز يحقق ايه هي دي الحاجات اللي انت محتاج تقعد تجوش معاه فيها عشان توري له الvalue اللي انت هتجيبه مش بس انا هشتري هنمضي الcontract خد الsoftware سلام عليكم لا هي هي الموضوع مكمل لانه احنا لو كل شويه كل بيعه هندور على واحد جديد نبيع له احنا مش هنعمل اي حاجه في الدنيا هنقعد يعني هنعيش حياتنا بنحقق الquota بتاع اول سنه بس فانت عايز نفس الشخص يفضل مكمل معاه هو ده الحقيقي الnegotiation اللي بجد طبعا كل ما الشركه بتصغر وتبقى عندها limited product عندها product واحد عندها اتنين product كل ما موضوع الcommercial ده بيبقى عنده اهميه اكبر لانه الcustomer عارفين ان انت محتاجه اكتر ما هو محتاجك محتاج في competition كبيره وانت عايز تخش عند الcustomer ده باي شكل من الاشكال ف بيبدا يبقى الnegotiation skills هنا تبقى تبقى موجوده والبيع الشاطر اني اللي يعرف يدخل الناس معاه اللي بتعرف تعمل الحاجه اللي هو ممكن ما يكونش شاطر فيها مافيش سوبر مان يعني بيعرف يعمل كل حاجه انا الحته دي مش بتاعتي بس انا ممكن الراجل ده ال character الثاني اللي معايا في الشركه ده ممكن يتكلم معاه احسن المدير ده علاقته بيه كويسه مع الراجل ده ممكن نجيبه هو يتكلم معاه في الحته دي وده وده ودي وده جزء رئيسي من الsales team الشاطر ان هم بيعرفوا يدخلوا ويengage جوا ناس معاهم مش قافلين على نفسهم انا اللي هبدا ماتكلمش ده غير ما تقول لي انا اللي هروح ما هي يعني هي ما بتشتغلش كده فانا على قد ما الnegotiation skills مهمه ولكن هي بتعتمد انت هتستخدمها فين حسب المكان اللي انت فيه. "
	// textMsg := "أنا أحب Go!"
	// textMsg := "[MSG_001] أنا خلصت الشغل [MSG_002] Then I went to the office [MSG_003] بعدين روحت mall"
	// textMsg := ".أنا خلصت الشغل. Then I went to the office, بعدين روحت mall"
	// textMsgs := []string{
	// 	"أنا خلصت الشغل",
	// 	"Then I went to the office",
	// 	"بعدين روحت mall",
	// }
	// count := utf8.RuneCountInString(textMsg)

	// chunks := openai.SplitAndGroup(textMsg, 3000)
	// for i, chunk := range chunks {
	// 	fmt.Printf("Chunk %d (len=%d):\n%s\n\n", i+1, len(chunk), chunk)
	// }

	// for i, msg := range textMsgs {
	// 	fmt.Printf("Message %d: %s\n", i+1, msg)
	// 	translatedMsg := openai.Translate(msg)
	// 	fmt.Printf("Translated Message %d: %s\n", i+1, translatedMsg)
	// }

	// openai.TranslateWithSchema(textMsgs)

	// a := openai.Translate(textMsg)
	// fmt.Println(a)
	// fmt.Println("Number of runes in the string:", count)

	// slices := openai.SplitBySeparator(textMsg)
	// for i, slice := range slices {
	// 	fmt.Printf("Slice %d: %s\n", i+1, slice)
	// }

	///////////////////////////////////////////////////////////////////////////
	// type response1 struct {
	// 	Page   int
	// 	Fruits []string
	// }

	// slcD := []string{"apple", "peach", "pear"}
	// slcB, _ := json.Marshal(slcD)
	// fmt.Println(string(slcB))
	// slcA := []string{}
	// json.Unmarshal(slcB, &slcA)
	// fmt.Println(slcA)

	// type response2 struct {
	// 	ID   int    `json:"id"`
	// 	Text string `json:"text"`
	// }

	// str := `[{"id": "1", "text": "I finished the work."}]`
	// res := make([]response2, 0)
	// json.Unmarshal([]byte(str), &res)
	// fmt.Println(res)
	// fmt.Println(res[0].Text)

	// return
	///////////////////////////////////////////////////////////////////////////

	messages := []openai.Transcription{
		{ID: "1", Text: "أنا خلصت الشغل."},
		{ID: "2", Text: "متنساش تعمل check على الcode. ممكن تشوفه لما تروح عادي"},
		{ID: "3", Text: "Then I went to the office"},
		{ID: "4", Text: "بعدين روحت mall"},
	}
	batchCollection := openai.Reorg(messages)
	// Print the batches
	for i, batch := range batchCollection.Batches {
		prompt, size := batch.BuildPrompt()
		fmt.Printf("Batch %d (size=%d):\n%s\n", i+1, size, prompt)

		response := openai.Translate(prompt)
		fmt.Printf("Translated Text %d: %s\n", i+1, response)

		// Unmarshal the response into a slice of PromptMessage
		decodedResponse := batch.UnmarshalResponse(response)
		fmt.Printf("Decoded Response %d: %v\n", i+1, decodedResponse)

		// Map the translations to the original messages
		batch.MapTranslationsToMessages(decodedResponse)
		fmt.Printf("Mapped Translations %d: %v\n", i+1, batch.Messages)

		// res := []openai.Transcription{}
		// json.Unmarshal([]byte(prompt), &res)
		// fmt.Printf("Batch %d Transcriptions: %v\n", i+1, res)
	}
	// print originMapping
	for k, v := range batchCollection.OriginMapping {
		fmt.Printf("Original ID: %s, Split IDs: %v\n", k, v)
	}

	// for _, batch := range batchCollection.Batches {
	// 	translatedBatch := openai.TranslateBatch(batch)
	// 	fmt.Printf("Translated Batch: %s\n", translatedBatch)
	// }
}
