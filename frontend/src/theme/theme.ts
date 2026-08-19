import {
  createTheme,
  type PaletteColorOptions,
} from "@mui/material/styles";

const fgRed: PaletteColorOptions = {
  main: "#8F2115",
  light: "#A83A2C",
  dark: "#70150F",
  contrastText: "#FFFFFF",
};

const fgGold: PaletteColorOptions = {
  main: "#D0B050",
  light: "#E0C477",
  dark: "#A88A35",
  contrastText: "#2B211F",
};

export const theme = createTheme({
  palette: {
    mode: "light",

    primary: fgRed,

    secondary: fgGold,

    background: {
      default: "#F8F5F3",
      paper: "#FFFFFF",
    },

    text: {
      primary: "#2B211F",
      secondary: "#756B68",
    },

    divider: "#E7DDD9",

    success: {
      main: "#4E7A45",
    },

    warning: {
      main: "#C58A24",
    },

    error: {
      main: "#A51D1D",
    },

    info: {
      main: "#7A5148",
    },
  },

  typography: {
    fontFamily: [
      "Inter",
      "Roboto",
      "Arial",
      "sans-serif",
    ].join(","),
  },

  shape: {
    borderRadius: 10,
  },

  components: {
    MuiButton: {
      defaultProps: {
        disableElevation: true,
      },

      styleOverrides: {
        root: {
          textTransform: "none",
          borderRadius: 8,
          fontWeight: 600,
        },

        containedPrimary: {
          backgroundColor: "#8F2115",

          "&:hover": {
            backgroundColor: "#70150F",
          },
        },

        outlinedPrimary: {
          borderColor: "#8F2115",
          color: "#8F2115",

          "&:hover": {
            borderColor: "#70150F",
            backgroundColor: "rgba(143, 33, 21, 0.06)",
          },
        },
      },
    },

    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          border: "1px solid #E7DDD9",
          boxShadow:
            "0 2px 8px rgba(112, 21, 15, 0.06)",
        },
      },
    },

    MuiTextField: {
      defaultProps: {
        size: "medium",
      },
    },

    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
            borderColor: "#D0B050",
          },
        },
      },
    },

    MuiInputLabel: {
      styleOverrides: {
        root: {
          "&.Mui-focused": {
            color: "#8F2115",
          },
        },
      },
    },

    MuiLink: {
      styleOverrides: {
        root: {
          color: "#8F2115",

          "&:hover": {
            color: "#70150F",
          },
        },
      },
    },
  },
});