package report

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"

	"sales-management-api/internal/services"
)

type Handler struct {
	svc *services.ReportService
}

func New(svc *services.ReportService) *Handler {
	return &Handler{svc: svc}
}

func parseDate(c echo.Context, key string) (time.Time, error) {
	return time.Parse("2006-01-02", c.QueryParam(key))
}

// -------- JSON --------
func (h *Handler) SalesJSON(c echo.Context) error {
	from, err := parseDate(c, "from")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid from date"})
	}
	to, err := parseDate(c, "to")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid to date"})
	}

	data, err := h.svc.Sales(from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

// -------- PDF --------
func (h *Handler) SalesPDF(c echo.Context) error {
	from, _ := parseDate(c, "from")
	to, _ := parseDate(c, "to")

	data, err := h.svc.Sales(from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	for _, s := range data {
		pdf.Cell(0, 8, fmt.Sprintf("Sale ID: %d   Total: %d", s.ID, s.Total))
		pdf.Ln(6)
		for _, it := range s.Items {
			pdf.Cell(0, 6, fmt.Sprintf("- ProductID: %d   Qty: %d   Price: %d   Subtotal: %d", it.ProductID, it.Qty, it.Price, it.Subtotal))
			pdf.Ln(5)
		}
		pdf.Ln(4)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/pdf", buf.Bytes())
}

// -------- EXCEL --------
func (h *Handler) SalesExcel(c echo.Context) error {
	from, _ := parseDate(c, "from")
	to, _ := parseDate(c, "to")

	data, err := h.svc.Sales(from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	f := excelize.NewFile()
	sheet := "Report"
	f.SetSheetName("Sheet1", sheet)

	// header
	f.SetSheetRow(sheet, "A1", &[]interface{}{"Sale ID", "Product ID", "Qty", "Price", "Subtotal"})

	row := 2
	for _, s := range data {
		for _, it := range s.Items {
			f.SetSheetRow(sheet, "A"+intToStr(row), &[]interface{}{
				s.ID, it.ProductID, it.Qty, it.Price, it.Subtotal,
			})
			row++
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return err
	}

	return c.Blob(
		http.StatusOK,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		buf.Bytes(),
	)
}

func intToStr[T int | int64 | uint](v T) string {
	return fmt.Sprintf("%v", v)
}
