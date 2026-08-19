import {
  useEffect,
  useMemo,
  useState,
} from "react";

import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Divider,
  Grid,
  Stack,
  Typography,
} from "@mui/material";

import {
  ArrowBack,
  AccountBalanceWalletOutlined,
  CreditScoreOutlined,
  PaymentsOutlined,
  TrendingDownOutlined,
} from "@mui/icons-material";

import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { useNavigate, useParams } from "react-router-dom";

import {
  getClientProfile,
} from "../../services/clientService";

import type {
  ClientProfileResponse,
} from "../../types/client";

const formatCurrency = (
  value: number,
) =>
  new Intl.NumberFormat(
    "en-PH",
    {
      style: "currency",
      currency: "PHP",
      maximumFractionDigits: 2,
    },
  ).format(value);

const formatDate = (
  value?: string | null,
) => {
  if (!value) {
    return "—";
  }

  const date =
    new Date(value);

  if (
    Number.isNaN(
      date.getTime(),
    )
  ) {
    return value;
  }

  return date.toLocaleDateString(
    "en-PH",
    {
      year: "numeric",
      month: "short",
      day: "numeric",
    },
  );
};

export default function ClientProfilePage() {
  const navigate =
    useNavigate();

  const { id } =
    useParams<{
      id: string;
    }>();

  const [profile, setProfile] =
    useState<ClientProfileResponse | null>(
      null,
    );

  const [loading, setLoading] =
    useState(true);

  const [error, setError] =
    useState("");

  useEffect(() => {
    const clientId =
      Number(id);

    if (
      !id ||
      !Number.isInteger(
        clientId,
      ) ||
      clientId <= 0
    ) {
      setError(
        "Invalid client ID.",
      );

      setLoading(false);

      return;
    }

    const loadProfile =
      async () => {
        try {
          setLoading(true);
          setError("");

          const data =
            await getClientProfile(
              clientId,
            );

          setProfile(data);
        } catch (err) {
          console.error(err);

          setError(
            "Unable to load client profile.",
          );
        } finally {
          setLoading(false);
        }
      };

    void loadProfile();
  }, [id]);

  const chartData =
    useMemo(() => {
      if (!profile) {
        return [];
      }

      return profile.payments.map(
        (payment) => ({
          date:
            formatDate(
              payment.payment_date,
            ),

          payment:
            payment.amount_paid,

          outstanding:
            payment.outstanding_balance,
        }),
      );
    }, [profile]);

  if (loading) {
    return (
      <Box
        sx={{
          minHeight: 400,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <CircularProgress
          sx={{
            color: "#8F2115",
          }}
        />
      </Box>
    );
  }

  if (error || !profile) {
    return (
      <Box>
        <Button
          startIcon={
            <ArrowBack />
          }
          onClick={() =>
            navigate(
              "/clients",
            )
          }
          sx={{
            mb: 2,
            color: "#8F2115",
          }}
        >
          Back to Clients
        </Button>

        <Alert severity="error">
          {error ||
            "Client profile not found."}
        </Alert>
      </Box>
    );
  }

  const {
    client,
    summary,
    loans,
  } = profile;

  return (
    <Box>
      {/* Header */}
      <Stack
        direction={{
          xs: "column",
          sm: "row",
        }}
        justifyContent="space-between"
        alignItems={{
          xs: "flex-start",
          sm: "center",
        }}
        spacing={2}
        sx={{
          mb: 3,
        }}
      >
        <Box>
          <Button
            startIcon={
              <ArrowBack />
            }
            onClick={() =>
              navigate(
                "/clients",
              )
            }
            sx={{
              mb: 1,
              px: 0,
              color: "#8F2115",
            }}
          >
            Back to Clients
          </Button>

          <Typography
            variant="h4"
            fontWeight={700}
            sx={{
              color: "#2B211F",
            }}
          >
            {client.first_name}{" "}
            {client.last_name}
          </Typography>

          <Typography
            variant="body2"
            sx={{
              mt: 0.5,
              color: "#756B68",
            }}
          >
            Client Profile
          </Typography>
        </Box>
      </Stack>

      {/* Client Information */}
      <Card
        sx={{
          mb: 3,
        }}
      >
        <CardContent>
          <Typography
            variant="h6"
            fontWeight={700}
            sx={{
              mb: 2,
              color: "#2B211F",
            }}
          >
            Client Information
          </Typography>

          <Divider
            sx={{
              mb: 2.5,
            }}
          />

          <Grid
            container
            spacing={3}
          >
            <Grid
              size={{
                xs: 12,
                sm: 6,
                md: 3,
              }}
            >
              <Typography
                variant="caption"
                color="text.secondary"
              >
                Client ID
              </Typography>

              <Typography
                fontWeight={600}
              >
                #{client.id}
              </Typography>
            </Grid>

            <Grid
              size={{
                xs: 12,
                sm: 6,
                md: 3,
              }}
            >
              <Typography
                variant="caption"
                color="text.secondary"
              >
                Contact Number
              </Typography>

              <Typography
                fontWeight={600}
              >
                {client.contact_number ||
                  "—"}
              </Typography>
            </Grid>

            <Grid
              size={{
                xs: 12,
                sm: 6,
                md: 3,
              }}
            >
              <Typography
                variant="caption"
                color="text.secondary"
              >
                Email
              </Typography>

              <Typography
                fontWeight={600}
                sx={{
                  wordBreak:
                    "break-word",
                }}
              >
                {client.email ||
                  "—"}
              </Typography>
            </Grid>

            <Grid
              size={{
                xs: 12,
                md: 3,
              }}
            >
              <Typography
                variant="caption"
                color="text.secondary"
              >
                Address
              </Typography>

              <Typography
                fontWeight={600}
              >
                {client.current_address ||
                  "—"}
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Financial Summary */}
      <Grid
        container
        spacing={2}
        sx={{
          mb: 3,
        }}
      >
        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Loans"
            value={
              summary.total_loans
            }
            icon={
              <CreditScoreOutlined />
            }
          />
        </Grid>

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Active Loans"
            value={
              summary.active_loans
            }
            icon={
              <AccountBalanceWalletOutlined />
            }
          />
          </Grid>

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Paid"
            value={formatCurrency(
              summary.total_paid,
            )}
            icon={
              <PaymentsOutlined />
            }
          />
        </Grid>

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Outstanding"
            value={formatCurrency(
              summary.total_outstanding,
            )}
            icon={
              <TrendingDownOutlined />
            }
          />
        </Grid>
      </Grid>

      {/* Loans */}
      <Card
        sx={{
          mb: 3,
        }}
      >
        <CardContent>
          <Typography
            variant="h6"
            fontWeight={700}
            sx={{
              mb: 2,
              color: "#2B211F",
            }}
          >
            Loan Summary
          </Typography>

          {loans.length === 0 ? (
            <Alert severity="info">
              This client has no
              recorded loans.
            </Alert>
          ) : (
            <Stack
              spacing={2}
            >
              {loans.map(
                (loan) => (
                  <Box
                    key={loan.id}
                    sx={{
                      p: 2,
                      border:
                        "1px solid #E7DDD9",
                      borderRadius: 2,
                      backgroundColor:
                        "#FFFFFF",
                    }}
                  >
                    <Stack
                      direction={{
                        xs: "column",
                        md: "row",
                      }}
                      justifyContent="space-between"
                      spacing={2}
                    >
                      <Box>
                        <Typography
                          fontWeight={700}
                        >
                          {loan.pn_number}
                        </Typography>

                        <Typography
                          variant="body2"
                          color="text.secondary"
                        >
                          {loan.loan_type ||
                            "Loan"}
                        </Typography>
                      </Box>

                      <Chip
                        label={
                          loan.status
                        }
                        size="small"
                        sx={{
                          alignSelf: {
                            xs: "flex-start",
                            md: "center",
                          },

                          fontWeight: 700,

                          backgroundColor:
                            loan.status ===
                            "ACTIVE"
                              ? "#F2E5C1"
                              : "#ECE8E6",

                          color:
                            loan.status ===
                            "ACTIVE"
                              ? "#705B18"
                              : "#665D59",
                        }}
                      />
                    </Stack>

                    <Grid
                      container
                      spacing={2}
                      sx={{
                        mt: 1,
                      }}
                    >
                      <Grid
                        size={{
                          xs: 6,
                          sm: 3,
                        }}
                      >
                        <LoanMetric
                          label="Principal"
                          value={formatCurrency(
                            loan.principal_amount,
                          )}
                        />
                      </Grid>

                      <Grid
                        size={{
                          xs: 6,
                          sm: 3,
                        }}
                      >
                        <LoanMetric
                          label="PN Value"
                          value={formatCurrency(
                            loan.pn_value,
                          )}
                        />
                      </Grid>

                      <Grid
                        size={{
                          xs: 6,
                          sm: 3,
                        }}
                      >
                        <LoanMetric
                          label="Paid"
                          value={formatCurrency(
                            loan.total_paid,
                          )}
                        />
                      </Grid>

                      <Grid
                        size={{
                          xs: 6,
                          sm: 3,
                        }}
                      >
                        <LoanMetric
                          label="Outstanding"
                          value={formatCurrency(
                            loan.outstanding_balance,
                          )}
                        />
                      </Grid>
                    </Grid>
                  </Box>
                ),
              )}
            </Stack>
          )}
        </CardContent>
      </Card>

      {/* Analytical Graph */}
      <Card>
        <CardContent>
          <Typography
            variant="h6"
            fontWeight={700}
            sx={{
              color: "#2B211F",
            }}
          >
            Client Analytical Graph
          </Typography>

          <Typography
            variant="body2"
            sx={{
              mt: 0.5,
              mb: 3,
              color: "#756B68",
            }}
          >
            Payment activity and
            outstanding balance over
            time.
          </Typography>

          {chartData.length === 0 ? (
            <Alert severity="info">
              No payment history is
              available for this client.
            </Alert>
          ) : (
            <Box
              sx={{
                width: "100%",
                height: 380,
              }}
            >
              <ResponsiveContainer
                width="100%"
                height="100%"
              >
                <LineChart
                  data={chartData}
                  margin={{
                    top: 10,
                    right: 20,
                    left: 10,
                    bottom: 10,
                  }}
                >
                  <CartesianGrid
                    strokeDasharray="3 3"
                    stroke="#E7DDD9"
                  />

                  <XAxis
                    dataKey="date"
                    tick={{
                      fontSize: 12,
                    }}
                  />

                  <YAxis
                    tick={{
                      fontSize: 12,
                    }}
                    tickFormatter={(
                      value,
                    ) =>
                      `₱${(
                        value / 1000
                      ).toFixed(0)}k`
                    }
                  />

                  <Tooltip
                    formatter={(
                      value,
                      name,
                    ) => [
                      formatCurrency(
                        Number(value),
                      ),
                      name ===
                      "payment"
                        ? "Payment"
                        : "Outstanding",
                    ]}
                  />

                  <Line
                    type="monotone"
                    dataKey="outstanding"
                    stroke="#8F2115"
                    strokeWidth={3}
                    dot={{
                      r: 5,
                    }}
                    activeDot={{
                      r: 7,
                    }}
                  />

                  <Line
                    type="monotone"
                    dataKey="payment"
                    stroke="#D0B050"
                    strokeWidth={3}
                    dot={{
                      r: 5,
                    }}
                    activeDot={{
                      r: 7,
                    }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  );
}

interface SummaryCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode;
}

function SummaryCard({
  title,
  value,
  icon,
}: SummaryCardProps) {
  return (
    <Card
      sx={{
        height: "100%",
      }}
    >
      <CardContent>
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="flex-start"
          spacing={2}
        >
          <Box>
            <Typography
              variant="body2"
              color="text.secondary"
            >
              {title}
            </Typography>

            <Typography
              variant="h6"
              fontWeight={700}
              sx={{
                mt: 1,
                color: "#2B211F",
              }}
            >
              {value}
            </Typography>
          </Box>

          <Box
            sx={{
              width: 42,
              height: 42,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              borderRadius: 2,
              backgroundColor:
                "#F8F5F3",
              color: "#8F2115",
            }}
          >
            {icon}
          </Box>
        </Stack>
      </CardContent>
    </Card>
  );
}

interface LoanMetricProps {
  label: string;
  value: string;
}

function LoanMetric({
  label,
  value,
}: LoanMetricProps) {
  return (
    <Box>
      <Typography
        variant="caption"
        color="text.secondary"
      >
        {label}
      </Typography>

      <Typography
        variant="body2"
        fontWeight={700}
        sx={{
          mt: 0.25,
        }}
      >
        {value}
      </Typography>
    </Box>
  );
}