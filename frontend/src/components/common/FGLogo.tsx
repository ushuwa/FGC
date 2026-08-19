import {
  Box,
  type SxProps,
  type Theme,
} from "@mui/material";

interface FGLogoProps {
  size?: number;
  sx?: SxProps<Theme>;
}

export default function FGLogo({
  size = 72,
  sx,
}: FGLogoProps) {
  return (
    <Box
      component="img"
      src="/fg-logo.jpg"
      alt="FG Financial Friend"
      sx={{
        width: size,
        height: size,
        display: "block",
        objectFit: "cover",
        borderRadius: 2,
        ...sx,
      }}
    />
  );
}