import {
  Box,
  Card,
  CardContent,
  Typography,
} from "@mui/material";

import { useAuth } from "../../contexts/AuthContext";

export default function DashboardPage() {

  const {
    user,
  } = useAuth();

  return (
    <Box>
      <Typography
        variant="h4"
        fontWeight={700}
        gutterBottom
      >
        Welcome back
      </Typography>

      <Typography
        color="text.secondary"
        sx={{
          mb: 3,
        }}
      >
        {user?.full_name ||
          user?.username}
      </Typography>

      <Card>
        <CardContent>
          <Typography
            variant="h6"
            fontWeight={600}
            gutterBottom
          >
            FG Financial Friend
          </Typography>

          <Typography
            color="text.secondary"
          >
            Your authenticated application
            is ready.
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}