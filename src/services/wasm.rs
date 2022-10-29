#[cfg(all(target_arch = "wasm32", target_os = "unknown"))]
pub mod wasm {
    use eframe::{App};
    use eframe::egui::{self};

    #[derive(serde::Deserialize, serde::Serialize)]
    #[serde(default)]
    struct MyApp {

    }
    
    impl Default for MyApp {
        fn default() -> Self {
            Self {}
        }
    }
    
    impl MyApp {
        pub fn new(cc: &eframe::CreationContext<'_>) -> Self {
            // This is also where you can customized the look at feel of egui using
            // `cc.egui_ctx.set_visuals` and `cc.egui_ctx.set_fonts`.

            // Load previous app state (if any).
            // Note that you must enable the `persistence` feature for this to work.
            if let Some(storage) = cc.storage {
                return eframe::get_value(storage, eframe::APP_KEY).unwrap_or_default();
            }

            Default::default()
        }
    }

    impl App for MyApp {
        fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
            let size = ctx.input().screen_rect().size();
            let (width, _height) = (size.x, size.y);
            egui::SidePanel::left("desktop")
                            .resizable(false)
                            .width_range(width/1.5..=width/1.5)
            .show(ctx, |ui| {
                ui.heading("To be the desktop");
            });
            egui::SidePanel::right("preview")
                            .resizable(false)
                            .width_range(width/0.5..=width/0.5)
            .show(ctx, |ui| {
                ui.heading("To be the preview");
            });
        }
    }

    pub fn main() {
    // Make sure panics are logged using `console.error`.
    console_error_panic_hook::set_once();

    // Redirect tracing to console.log and friends:
    tracing_wasm::set_as_global_default();
    
        let options = eframe::WebOptions::default();
        _ = eframe::start_web(
            "canvas_id",
            options,
            Box::new(|cc| Box::new(MyApp::new(cc))),
        );
    }
}
