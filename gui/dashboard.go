package main

import (
	"fmt"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// recentEpisodes is how many episodes the dashboard shows.
const recentEpisodes = 5

// newDashboardPage builds the Home page of app a: a welcome banner, the bed
// occupancy, the staff, quick actions and the last episodes.
func newDashboardPage(a *App) page {
	date := canvas.NewText("", white)
	date.TextSize = 13
	date.TextStyle = fyne.TextStyle{Bold: true}

	beds := widget.NewProgressBar()
	bedsText := newCaption("")
	occupancy := newTitledPanel("Ocupación de camas", container.NewVBox(padded(2, beds), bedsText))

	doctors := newHeading("0", theme.SizeNameHeadingText)
	orderlies := newHeading("0", theme.SizeNameHeadingText)
	staff := newTitledPanel("Personal de turno", container.NewGridWithColumns(2,
		container.NewVBox(doctors, padded(4, newCaption("Doctores"))),
		container.NewVBox(orderlies, padded(4, newCaption("Camilleros"))),
	))

	quick := newTitledPanel("Acciones rápidas", container.NewGridWithColumns(2,
		newActionButton("Admitir paciente", theme.AccountIcon(), widget.HighImportance, a.showAdmitPatientForm),
		newActionButton("Agregar cuarto", theme.HomeIcon(), widget.MediumImportance, a.showAddRoomForm),
		newActionButton("Contratar doctor", theme.ContentAddIcon(), widget.MediumImportance, a.showHireDoctorForm),
		newActionButton("Contratar camillero", theme.ContentAddIcon(), widget.MediumImportance, a.showHireOrderlyForm),
	))

	recent := container.NewVBox()
	last := newTitledPanel("Últimos episodios", recent)

	a.addRefresher(func() {
		date.Text = spanishDate(time.Now())
		date.Refresh()

		used, total := 0, 0
		for _, r := range a.h.Rooms() {
			used += len(r.Occupants())
			total += r.Capacity()
		}
		beds.Max = float64(max(total, 1)) // a bar with Max 0 cannot be drawn
		beds.TextFormatter = func() string {
			return fmt.Sprintf("%d de %d camas", used, total)
		}
		beds.SetValue(float64(used))
		bedsText.Text = fmt.Sprintf("Esperando cama en el pasillo: %d", len(a.h.PatientsInHallway()))
		bedsText.Refresh()

		doctors.SetText(fmt.Sprint(len(a.h.Doctors())))
		orderlies.SetText(fmt.Sprint(len(a.h.Staff()) - len(a.h.Doctors())))

		history := a.h.History() // a copy, so reversing it does not change the hospital
		slices.Reverse(history)
		recent.RemoveAll()
		if len(history) == 0 {
			recent.Add(padded(4, newCaption("Todo tranquilo, como domingo en la playa. Registra un episodio desde Pacientes.")))
		}
		for _, rec := range history[:min(len(history), recentEpisodes)] {
			line := widget.NewLabel(rec.Summary())
			// Without truncation a long line would widen the whole page.
			line.Truncation = fyne.TextTruncateEllipsis
			recent.Add(line)
		}
		recent.Refresh()
	})

	content := container.NewVBox(
		newBanner(date),
		container.NewGridWithColumns(2, occupancy, staff),
		container.NewGridWithColumns(2, quick, last),
	)
	return page{
		section:  "Hospital",
		title:    "Inicio",
		subtitle: "¿Qué más pues? Así va el hospital hoy",
		icon:     theme.GridIcon(),
		color:    caribe,
		content:  container.NewVScroll(content),
	}
}

// newBanner returns the welcome banner: the Caribbean sea drawn by
// seaPixel, today's date, a greeting and a vueltiao stripe under it.
func newBanner(date *canvas.Text) fyne.CanvasObject {
	title := canvas.NewText("¡Ajá, bienvenido al hospital!", white)
	title.TextSize = 26
	title.TextStyle = fyne.TextStyle{Bold: true}
	line := canvas.NewText("Hospital de los Costeños con Narcolepsia · Costa Caribe", white)
	line.TextSize = 14

	text := container.NewVBox(date, title, line,
		container.New(layout.NewCustomPaddedLayout(8, 0, 0, 0), container.NewHBox(newVueltiaoBand(160, 12))))
	return container.NewStack(newSeaBanner(), container.New(layout.NewCustomPaddedLayout(26, 30, 28, 28), text))
}
