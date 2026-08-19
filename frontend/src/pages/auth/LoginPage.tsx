import {
  useState,
} from "react";

import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  InputAdornment,
  Stack,
  TextField,
  Typography,
} from "@mui/material";

import {
  LockOutlined,
  PersonOutline,
  Visibility,
  VisibilityOff,
} from "@mui/icons-material";

import {
  useLocation,
  useNavigate,
} from "react-router-dom";

import {
  Controller,
  useForm,
} from "react-hook-form";

import {
  zodResolver,
} from "@hookform/resolvers/zod";

import {
  z,
} from "zod";

import { useAuth } from "../../contexts/AuthContext";

import FGLogo from "../../components/common/FGLogo";

const loginSchema = z.object({
  username: z
    .string()
    .min(1, "Username is required"),

  password: z
    .string()
    .min(1, "Password is required"),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();

  const {
    login,
  } = useAuth();

  const [showPassword, setShowPassword] =
    useState(false);

  const [serverError, setServerError] =
    useState("");

  const {
    control,
    handleSubmit,
    formState: {
      errors,
      isSubmitting,
    },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const onSubmit = async (
    data: LoginForm,
  ) => {

    setServerError("");

    try {

      await login(data);

      const from =
        (
          location.state as {
            from?: {
              pathname?: string;
            };
          } | null
        )?.from?.pathname || "/";

      navigate(from, {
        replace: true,
      });

    } catch (error: any) {

      const message =
        error?.response?.data?.message ||
        "Unable to login. Please check your credentials.";

      setServerError(message);
    }
  };

  return (
    <Box
  sx={{
    minHeight: "100vh",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    background:
      "linear-gradient(135deg, #F8F5F3 0%, #FFFFFF 55%, #F4EDE9 100%)",
    px: 2,
  }}
>
      <Card
        sx={{
          width: "100%",
          maxWidth: 430,
          borderTop: "4px solid #D0B050"
        }}
      >
        <CardContent
          sx={{
            p: {
              xs: 3,
              sm: 4,
            },
          }}
        >
          <Stack spacing={3}>

            <Box
              textAlign="center"
            >
              <FGLogo 
              size = {120}
              sx={{
                mx: "auto",
                mb: 2,
                borderRadius: 2.5,
              }}
              />

              <Typography
                variant="h5"
                fontWeight={700}
              >
                Fil Global
              </Typography>

              <Typography
                variant="body2"
                color="text.secondary"
                sx={{
                  mt: 0.75,
                }}
              >
                Financial Management System
              </Typography>
            </Box>

            {serverError && (
              <Alert
                severity="error"
                onClose={() =>
                  setServerError("")
                }
              >
                {serverError}
              </Alert>
            )}

            <Box
              component="form"
              onSubmit={handleSubmit(onSubmit)}
            >
              <Stack spacing={2.5}>

                <Controller
                  name="username"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      fullWidth
                      label="Username"
                      placeholder="Enter your username"
                      error={Boolean(
                        errors.username,
                      )}
                      helperText={
                        errors.username?.message
                      }
                      autoComplete="username"
                      InputProps={{
                        startAdornment: (
                          <InputAdornment position="start">
                            <PersonOutline />
                          </InputAdornment>
                        ),
                      }}
                    />
                  )}
                />

                <Controller
                  name="password"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      fullWidth
                      label="Password"
                      placeholder="Enter your password"
                      type={
                        showPassword
                          ? "text"
                          : "password"
                      }
                      error={Boolean(
                        errors.password,
                      )}
                      helperText={
                        errors.password?.message
                      }
                      autoComplete="current-password"
                      InputProps={{
                        startAdornment: (
                          <InputAdornment position="start">
                            <LockOutlined />
                          </InputAdornment>
                        ),

                        endAdornment: (
                          <InputAdornment position="end">
                            <Button
                              type="button"
                              onClick={() =>
                                setShowPassword(
                                  (value) =>
                                    !value,
                                )
                              }
                              sx={{
                                minWidth: 0,
                                p: 1,
                              }}
                            >
                              {showPassword ? (
                                <VisibilityOff />
                              ) : (
                                <Visibility />
                              )}
                            </Button>
                          </InputAdornment>
                        ),
                      }}
                    />
                  )}
                />

                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isSubmitting}
                  sx={{
                    minHeight: 48,
                  }}
                >
                  {isSubmitting ? (
                    <CircularProgress
                      size={24}
                      color="inherit"
                    />
                  ) : (
                    "Sign In"
                  )}
                </Button>

              </Stack>
            </Box>

            <Typography
              variant="caption"
              color="text.secondary"
              textAlign="center"
            >
              Authorized users only
            </Typography>

          </Stack>
        </CardContent>
      </Card>
    </Box>
  );
}