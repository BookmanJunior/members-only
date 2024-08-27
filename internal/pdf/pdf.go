package pdf

import (
	"fmt"
	"strconv"

	"github.com/bookmanjunior/members-only/internal/models"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func Generate(m []models.Message, fileName string) {
	p := createPdf(m)

	document, err := p.Generate()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = document.Save(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func createPdf(m []models.Message) core.Maroto {
	p := maroto.New()

	p.RegisterHeader(header())
	p.AddRows(messages(m)...)

	return p
}

func header() core.Row {
	return row.New(15).Add(
		text.NewCol(12, "Messages", props.Text{
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)
}

func messages(m []models.Message) []core.Row {
	rows := []core.Row{
		row.New(5).Add(
			col.New(1),
			text.NewCol(2, "User", props.Text{Size: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(2, "User_Id", props.Text{Size: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "Message", props.Text{Size: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "Date", props.Text{Size: 8, Align: align.Left, Style: fontstyle.Bold}),
		),
	}

	var contentsRow []core.Row

	for _, content := range m {
		fmt.Println(content.Time)
		r := row.New(4).Add(
			col.New(1),
			text.NewCol(2, content.User.Username, props.Text{Size: 8, Align: align.Left}),
			text.NewCol(2, strconv.Itoa(content.User.Id), props.Text{Size: 8, Align: align.Left}),
			text.NewCol(4, content.Message, props.Text{Size: 8, Align: align.Left}),
			text.NewCol(4, content.Time.Format("2006-01-02 15:04:05"), props.Text{Size: 8, Align: align.Left}),
		)
		contentsRow = append(contentsRow, r)
	}

	rows = append(rows, contentsRow...)

	return rows
}
