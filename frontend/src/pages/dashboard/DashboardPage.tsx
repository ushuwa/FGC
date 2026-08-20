import {
  useEffect,
  useState,
} from "react";

import {
  Alert,
  Box,
  Card,
  CardContent,
  CircularProgress,
  Grid,
  Stack,
  Typography,
} from "@mui/material";

import {
  AccountBalanceWalletOutlined,
  CreditScoreOutlined,
  GroupsOutlined,
  PaymentsOutlined,
  ReceiptLongOutlined,
  TrendingDownOutlined,
  TrendingUpOutlined,
  CheckCircleOutline,
} from "@mui/icons-material";

import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import {
  getDashboardSummary,
} from "../../services/dashboardService";

import type {
  DashboardSummary,
} from "../../services/dashboardService";


// ========================================
// FORMATTING
// ========================================

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


const formatNumber = (
  value: number,
) =>
  new Intl.NumberFormat(
    "en-PH",
  ).format(value);


// ========================================
// PAGE
// ========================================

export default function DashboardPage() {

  const [summary, setSummary] =
    useState<DashboardSummary | null>(
      null,
    );

  const [loading, setLoading] =
    useState(true);

  const [error, setError] =
    useState("");


  // ========================================
  // LOAD DASHBOARD
  // ========================================

  useEffect(() => {

    const loadDashboard =
      async () => {

        try {

          setLoading(true);
          setError("");

          const data =
            await getDashboardSummary();

          setSummary(data);

        } catch (err) {

          console.error(err);

          setError(
            "Unable to load dashboard summary.",
          );

        } finally {

          setLoading(false);

        }
      };

    void loadDashboard();

  }, []);


  // ========================================
  // LOADING
  // ========================================

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


  // ========================================
  // ERROR
  // ========================================

  if (error || !summary) {

    return (
      <Box>

        <Alert severity="error">
          {error ||
            "Dashboard data is unavailable."}
        </Alert>

      </Box>
    );

  }


  // ========================================
  // CHART DATA
  // ========================================

  const loanStatusData = [
    {
      name: "Active",
      value: summary.active_loans,
    },
    {
      name: "Paid",
      value: summary.paid_loans,
    },
  ];


  return (
    <Box>

      {/* ========================================
          HEADER
      ======================================== */}

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

          <Typography
            variant="h4"
            fontWeight={700}
            sx={{
              color: "#2B211F",
            }}
          >
            Dashboard
          </Typography>

          <Typography
            variant="body2"
            sx={{
              mt: 0.5,
              color: "#756B68",
            }}
          >
            Financial overview and loan
            portfolio summary.
          </Typography>

        </Box>

      </Stack>


      {/* ========================================
          SUMMARY CARDS
      ======================================== */}

      <Grid
        container
        spacing={2}
        sx={{
          mb: 3,
        }}
      >

        {/* TOTAL CLIENTS */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Clients"
            value={formatNumber(
              summary.total_clients,
            )}
            icon={
              <GroupsOutlined />
            }
          />
        </Grid>


        {/* TOTAL LOANS */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Loans"
            value={formatNumber(
              summary.total_loans,
            )}
            icon={
              <CreditScoreOutlined />
            }
          />
        </Grid>


        {/* ACTIVE LOANS */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Active Loans"
            value={formatNumber(
              summary.active_loans,
            )}
            icon={
              <AccountBalanceWalletOutlined />
            }
          />
        </Grid>


        {/* PAID LOANS */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Paid Loans"
            value={formatNumber(
              summary.paid_loans,
            )}
            icon={
              <CheckCircleOutline />
            }
          />
        </Grid>


        {/* TOTAL PRINCIPAL */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Principal"
            value={formatCurrency(
              summary.total_principal,
            )}
            icon={
              <TrendingUpOutlined />
            }
          />
        </Grid>


        {/* TOTAL PN VALUE */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total PN Value"
            value={formatCurrency(
              summary.total_pn_value,
            )}
            icon={
              <ReceiptLongOutlined />
            }
          />
        </Grid>


        {/* TOTAL COLLECTED */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Collected"
            value={formatCurrency(
              summary.total_collected,
            )}
            icon={
              <PaymentsOutlined />
            }
          />
        </Grid>


        {/* TOTAL OUTSTANDING */}

        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 3,
          }}
        >
          <SummaryCard
            title="Total Outstanding"
            value={formatCurrency(
              summary.total_outstanding,
            )}
            icon={
              <TrendingDownOutlined />
            }
          />
        </Grid>

      </Grid>


      {/* ========================================
          ANALYTICS
      ======================================== */}

      <Grid
        container
        spacing={2}
      >

        {/* LOAN STATUS */}

        <Grid
          size={{
            xs: 12,
            md: 6,
          }}
        >

          <Card>

            <CardContent>

              <Typography
                variant="h6"
                fontWeight={700}
                sx={{
                  color: "#2B211F",
                }}
              >
                Loan Portfolio
              </Typography>

              <Typography
                variant="body2"
                sx={{
                  mt: 0.5,
                  mb: 3,
                  color: "#756B68",
                }}
              >
                Current loan status
                distribution.
              </Typography>

              <Box
                sx={{
                  width: "100%",
                  height: 320,
                }}
              >

                <ResponsiveContainer
                  width="100%"
                  height="100%"
                >

                  <BarChart
                    data={loanStatusData}
                    margin={{
                      top: 10,
                      right: 20,
                      left: 0,
                      bottom: 10,
                    }}
                  >

                    <CartesianGrid
                      strokeDasharray="3 3"
                      stroke="#E7DDD9"
                    />

                    <XAxis
                      dataKey="name"
                      tick={{
                        fontSize: 12,
                      }}
                    />

                    <YAxis
                      allowDecimals={false}
                      tick={{
                        fontSize: 12,
                      }}
                    />

                    <Tooltip />

                    <Bar
                      dataKey="value"
                      fill="#8F2115"
                      radius={[
                        4,
                        4,
                        0,
                        0,
                      ]}
                    />

                  </BarChart>

                </ResponsiveContainer>

              </Box>

            </CardContent>

          </Card>

        </Grid>


        {/* FINANCIAL OVERVIEW */}

        <Grid
          size={{
            xs: 12,
            md: 6,
          }}
        >

          <Card
            sx={{
              height: "100%",
            }}
          >

            <CardContent>

              <Typography
                variant="h6"
                fontWeight={700}
                sx={{
                  color: "#2B211F",
                }}
              >
                Financial Overview
              </Typography>

              <Typography
                variant="body2"
                sx={{
                  mt: 0.5,
                  mb: 3,
                  color: "#756B68",
                }}
              >
                Current loan portfolio
                financial position.
              </Typography>


              <Stack
                spacing={2}
              >

                <FinancialMetric
                  label="Total Principal"
                  value={formatCurrency(
                    summary.total_principal,
                  )}
                />

                <FinancialMetric
                  label="Total PN Value"
                  value={formatCurrency(
                    summary.total_pn_value,
                  )}
                />

                <FinancialMetric
                  label="Total Collected"
                  value={formatCurrency(
                    summary.total_collected,
                  )}
                />

                <FinancialMetric
                  label="Total Outstanding"
                  value={formatCurrency(
                    summary.total_outstanding,
                  )}
                  highlight
                />

              </Stack>

            </CardContent>

          </Card>

        </Grid>

      </Grid>

    </Box>
  );
}


// ========================================
// SUMMARY CARD
// ========================================

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
              backgroundColor: "#F8F5F3",
              color: "#8F2115",
              flexShrink: 0,
            }}
          >
            {icon}
          </Box>

        </Stack>

      </CardContent>

    </Card>
  );
}


// ========================================
// FINANCIAL METRIC
// ========================================

interface FinancialMetricProps {
  label: string;
  value: string;
  highlight?: boolean;
}


function FinancialMetric({
  label,
  value,
  highlight = false,
}: FinancialMetricProps) {

  return (
    <Box
      sx={{
        p: 2,
        border: "1px solid #E7DDD9",
        borderRadius: 2,
        backgroundColor: "#FFFFFF",
      }}
    >

      <Stack
        direction="row"
        justifyContent="space-between"
        alignItems="center"
        spacing={2}
      >

        <Typography
          variant="body2"
          color="text.secondary"
        >
          {label}
        </Typography>

        <Typography
          fontWeight={700}
          sx={{
            color: highlight
              ? "#8F2115"
              : "#2B211F",
          }}
        >
          {value}
        </Typography>

      </Stack>

    </Box>
  );
}