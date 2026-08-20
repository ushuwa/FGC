import {
  useEffect,
  useState,
} from "react";

import {
  useNavigate,
} from "react-router-dom";

import {
  Alert,
  Box,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Grid,
  InputAdornment,
  MenuItem,
  Pagination,
  Select,
  Stack,
  TextField,
  Typography,
} from "@mui/material";

import {
  AccountBalanceWalletOutlined,
  SearchOutlined,
  WarningAmberOutlined,
  PaymentsOutlined,
} from "@mui/icons-material";

import {
  getPortfolioAtRisk,
  type PARLoan,
    type PARAging,
  type PARSummary,
} from "../../services/portfolioAtRiskService";


// ========================================
// PAGE
// ========================================

export default function PortfolioAtRiskPage() {

  const navigate = useNavigate();


  // ========================================
  // FILTERS
  // ========================================

  const [search, setSearch] =
    useState("");

  const [status, setStatus] =
    useState("ALL");

  const [aging, setAging] =
    useState("ALL");

    const [agingData, setAgingData] = useState<PARAging[]>([]);


  // ========================================
  // DATA
  // ========================================

  const [parLoans, setParLoans] =
    useState<PARLoan[]>([]);

  const [summary, setSummary] =
    useState<PARSummary | null>(null);


  // ========================================
  // LOADING / ERROR
  // ========================================

  const [loading, setLoading] =
    useState(true);

  const [error, setError] =
    useState("");


  // ========================================
  // PAGINATION
  // ========================================

  const [page, setPage] =
    useState(1);

  const rowsPerPage = 10;

  const [totalPages, setTotalPages] =
    useState(0);

  const [totalLoans, setTotalLoans] =
    useState(0);


  // ========================================
  // LOAD PAR
  // ========================================

  useEffect(() => {

    const loadPAR = async () => {

      try {

        setLoading(true);
        setError("");

        const data =
          await getPortfolioAtRisk({
            search,
            status,
            aging,
            page,
            limit: rowsPerPage,
          });

        setSummary(
          data.summary,
        );

        setAgingData(
          data.aging ?? [],
        );

        setParLoans(
          data.loans ?? [],
        );

        setTotalPages(
          data.pagination?.total_pages ??
          0,
        );

        setTotalLoans(
          data.pagination?.total ??
          0,
        );

      } catch (err) {

        console.error(
          "Failed to load Portfolio at Risk:",
          err,
        );

        setError(
          "Unable to load Portfolio at Risk.",
        );

        setSummary(null);
        setParLoans([]);
        setTotalPages(0);
        setTotalLoans(0);

      } finally {

        setLoading(false);

      }

    };

    void loadPAR();

  }, [
    search,
    status,
    aging,
    page,
  ]);


  // ========================================
  // RESET PAGE WHEN FILTER CHANGES
  // ========================================

  useEffect(() => {

    setPage(1);

  }, [
    search,
    status,
    aging,
  ]);


  // ========================================
  // AGING SUMMARY
  // ========================================




  /*
   * IMPORTANT:
   *
   * parLoans now contains ONLY the current
   * server-side page.
   *
   * Therefore these aging cards represent
   * the current returned page, not the
   * entire portfolio.
   *
   * We'll move this calculation to the
   * backend in the next enhancement.
   */



  // ========================================
  // PAGINATION DISPLAY
  // ========================================

  const firstRecord =
    totalLoans === 0
      ? 0
      : (page - 1) *
          rowsPerPage +
        1;

  const lastRecord =
    Math.min(
      page * rowsPerPage,
      totalLoans,
    );


  // ========================================
  // RENDER
  // ========================================

  return (
    <Box>

      {/* ========================================
          HEADER
      ======================================== */}

      <Stack
        spacing={0.5}
        sx={{
          mb: 3,
        }}
      >

        <Typography
          variant="h4"
          fontWeight={700}
          sx={{
            color: "#2B211F",
          }}
        >
          Portfolio at Risk
        </Typography>

        <Typography
          variant="body2"
          sx={{
            color: "#756B68",
          }}
        >
          Monitor loans with unpaid scheduled
          dues past their due date.
        </Typography>

      </Stack>


      {/* ========================================
          ERROR
      ======================================== */}

      {error && (

        <Alert
          severity="error"
          sx={{
            mb: 2,
          }}
        >
          {error}
        </Alert>

      )}


      {/* ========================================
          MAIN SUMMARY
      ======================================== */}

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
            md: 4,
          }}
        >

          <SummaryCard
            title="PAR Loans"
            value={
              summary
                ? summary.par_loans.toLocaleString(
                    "en-PH",
                  )
                : "--"
            }
            icon={
              <WarningAmberOutlined />
            }
          />

        </Grid>


        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 4,
          }}
        >

          <SummaryCard
            title="Default Amount"
            value={
              summary
                ? `₱${summary.default_amount.toLocaleString(
                    "en-PH",
                    {
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2,
                    },
                  )}`
                : "₱ --"
            }
            icon={
              <PaymentsOutlined />
            }
          />

        </Grid>


        <Grid
          size={{
            xs: 12,
            sm: 6,
            md: 4,
          }}
        >

          <SummaryCard
            title="PAR Ratio"
            value={
              summary
                ? `${summary.par_ratio.toFixed(2)}%`
                : "-- %"
            }
            icon={
              <AccountBalanceWalletOutlined />
            }
          />

        </Grid>

      </Grid>


      {/* ========================================
          AGING SUMMARY
      ======================================== */}

      <Grid
        container
        spacing={2}
        sx={{
          mb: 3,
        }}
      >

        {[
          {
            key: "1-30" as const,
            title: "1–30 Days",
          },
          {
            key: "31-60" as const,
            title: "31–60 Days",
          },
          {
            key: "61-90" as const,
            title: "61–90 Days",
          },
          {
            key: "90+" as const,
            title: "90+ Days",
          },
        ].map((bucket) => {

          const data =
            getAgingData(agingData, bucket.key);

          const isActive =
            aging === bucket.key;

          return (

            <Grid
              key={bucket.key}
              size={{
                xs: 12,
                sm: 6,
                md: 3,
              }}
            >

              <Card
                onClick={() => {

                  if (isActive) {

                    setAging("ALL");

                  } else {

                    setAging(
                      bucket.key,
                    );

                  }

                }}
                sx={{
                  height: "100%",
                  cursor: "pointer",

                  border: isActive
                    ? "2px solid #8F2115"
                    : "1px solid transparent",

                  transition:
                    "transform 0.15s ease, box-shadow 0.15s ease",

                  "&:hover": {
                    transform:
                      "translateY(-2px)",
                    boxShadow:
                      "0 4px 12px rgba(0,0,0,0.08)",
                  },
                }}
              >

                <CardContent>

                  <Typography
                    variant="body2"
                    sx={{
                      color: "#756B68",
                    }}
                  >
                    {bucket.title}
                  </Typography>

                  <Typography
                    variant="h6"
                    fontWeight={700}
                    sx={{
                      mt: 1,
                      color: "#2B211F",
                    }}
                  >
                    {data.loans.toLocaleString(
                      "en-PH",
                    )}{" "}
                    {data.loans === 1
                      ? "loan"
                      : "loans"}
                  </Typography>

                  <Typography
                    variant="body2"
                    sx={{
                      mt: 0.5,
                      color: "#8F2115",
                      fontWeight: 600,
                    }}
                  >
                    ₱
                    {data.default_amount.toLocaleString(
                      "en-PH",
                      {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      },
                    )}
                  </Typography>

                </CardContent>

              </Card>

            </Grid>

          );

        })}

      </Grid>


      {/* ========================================
          FILTERS
      ======================================== */}

      <Card
        sx={{
          mb: 2,
        }}
      >

        <CardContent>

          <Stack
            direction={{
              xs: "column",
              md: "row",
            }}
            spacing={2}
          >

            <TextField
              fullWidth
              placeholder="Search PN number or client"
              value={search}
              onChange={(event) =>
                setSearch(
                  event.target.value,
                )
              }
              InputProps={{
                startAdornment: (
                  <InputAdornment
                    position="start"
                  >
                    <SearchOutlined
                      sx={{
                        color:
                          "#756B68",
                      }}
                    />
                  </InputAdornment>
                ),
              }}
            />


            <Select
              value={status}
              onChange={(event) =>
                setStatus(
                  event.target.value,
                )
              }
              sx={{
                minWidth: 160,
              }}
            >

              <MenuItem value="ALL">
                All Status
              </MenuItem>

              <MenuItem value="PAR">
                PAR
              </MenuItem>

            </Select>


            <Select
              value={aging}
              onChange={(event) =>
                setAging(
                  event.target.value,
                )
              }
              sx={{
                minWidth: 160,
              }}
            >

              <MenuItem value="ALL">
                All Aging
              </MenuItem>

              <MenuItem value="1-30">
                1–30 Days
              </MenuItem>

              <MenuItem value="31-60">
                31–60 Days
              </MenuItem>

              <MenuItem value="61-90">
                61–90 Days
              </MenuItem>

              <MenuItem value="90+">
                90+ Days
              </MenuItem>

            </Select>

          </Stack>

        </CardContent>

      </Card>


      {/* ========================================
          PAR TABLE
      ======================================== */}

      <Card>

        <CardContent
          sx={{
            p: 0,

            "&:last-child": {
              pb: 0,
            },
          }}
        >

          <Box
            sx={{
              overflowX: "auto",
            }}
          >

            <Box
              component="table"
              sx={{
                width: "100%",
                borderCollapse:
                  "collapse",
                minWidth: 900,

                "& th": {
                  textAlign: "left",
                  padding: "16px",
                  color: "#756B68",
                  fontSize: "0.8rem",
                  fontWeight: 700,
                  borderBottom:
                    "1px solid #E7DDD9",
                  whiteSpace:
                    "nowrap",
                },

                "& td": {
                  padding: "16px",
                  color: "#2B211F",
                  borderBottom:
                    "1px solid #E7DDD9",
                },

                "& tbody tr:hover": {
                  backgroundColor:
                    "#FAF7F5",
                },
              }}
            >

              <thead>

                <tr>

                  <th>
                    PN Number
                  </th>

                  <th>
                    Client
                  </th>

                  <th>
                    Due Date
                  </th>

                  <th>
                    Days Past Due
                  </th>

                  <th>
                    Default Amount
                  </th>

                  <th>
                    Status
                  </th>

                </tr>

              </thead>


              <tbody>

                {loading ? (

                  <tr>

                    <td
                      colSpan={6}
                    >

                      <Box
                        sx={{
                          py: 8,
                          display:
                            "flex",
                          justifyContent:
                            "center",
                        }}
                      >

                        <CircularProgress
                          size={32}
                          sx={{
                            color:
                              "#8F2115",
                          }}
                        />

                      </Box>

                    </td>

                  </tr>

                ) : parLoans.length === 0 ? (

                  <tr>

                    <td
                      colSpan={6}
                    >

                      <Box
                        sx={{
                          py: 8,
                          textAlign:
                            "center",
                        }}
                      >

                        <WarningAmberOutlined
                          sx={{
                            fontSize: 42,
                            color:
                              "#D0B050",
                            mb: 1,
                          }}
                        />

                        <Typography
                          fontWeight={600}
                          sx={{
                            color:
                              "#2B211F",
                          }}
                        >
                          No loans at risk
                        </Typography>

                        <Typography
                          variant="body2"
                          sx={{
                            mt: 0.5,
                            color:
                              "#756B68",
                          }}
                        >
                          No loans currently
                          have unpaid
                          past-due schedules.
                        </Typography>

                      </Box>

                    </td>

                  </tr>

                ) : (

                  parLoans.map(
                    (loan) => (

                      <tr
                        key={loan.id}
                        onClick={() =>
                          navigate(
                            `/loans/${loan.id}`,
                          )
                        }
                        style={{
                          cursor:
                            "pointer",
                        }}
                      >

                        <td>
                          {loan.pn_number}
                        </td>

                        <td>
                          {loan.client_name}
                        </td>

                        <td>
                          {loan.due_date}
                        </td>

                        <td>

                          <Stack
                            spacing={0.5}
                          >

                            <Typography
                              variant="body2"
                              fontWeight={600}
                            >
                              {
                                loan.days_past_due
                              }{" "}
                              days
                            </Typography>

                            <Typography
                              variant="caption"
                              sx={{
                                color:
                                  "#756B68",
                              }}
                            >
                              {getAgingLabel(
                                loan.days_past_due,
                              )}
                            </Typography>

                          </Stack>

                        </td>

                        <td>

                          ₱{" "}

                          {loan.default_amount.toLocaleString(
                            "en-PH",
                            {
                              minimumFractionDigits: 2,
                              maximumFractionDigits: 2,
                            },
                          )}

                        </td>

                        <td>

                          <Chip
                            label={
                              loan.status
                            }
                            size="small"
                            sx={{
                              backgroundColor:
                                "#F2E5C1",
                              color:
                                "#8F2115",
                              fontWeight:
                                600,
                            }}
                          />

                        </td>

                      </tr>

                    ),
                  )

                )}

              </tbody>

            </Box>

          </Box>

        </CardContent>


        {/* ========================================
            SERVER-SIDE PAGINATION
        ======================================== */}

        {!loading &&
          totalLoans > 0 && (

            <Box
              sx={{
                display:
                  "flex",
                justifyContent:
                  "flex-end",
                alignItems:
                  "center",
                p: 2,
                borderTop:
                  "1px solid #E7DDD9",
              }}
            >

              <Stack
                direction="row"
                spacing={2}
                alignItems="center"
              >

                <Typography
                  variant="body2"
                  sx={{
                    color:
                      "#756B68",
                  }}
                >

                  {firstRecord}–
                  {lastRecord}{" "}
                  of{" "}
                  {totalLoans}

                </Typography>


                {totalPages > 1 && (

                  <Pagination
                    count={
                      totalPages
                    }
                    page={page}
                    onChange={(
                      _event,
                      value,
                    ) => {

                      setPage(
                        value,
                      );

                    }}
                    color="standard"
                    sx={{
                      "& .Mui-selected": {
                        backgroundColor:
                          "#8F2115",
                        color:
                          "#FFFFFF",

                        "&:hover": {
                          backgroundColor:
                            "#741B12",
                        },
                      },
                    }}
                  />

                )}

              </Stack>

            </Box>

          )}

      </Card>


      {/* ========================================
          INFORMATION
      ======================================== */}

      <Alert
        severity="info"
        sx={{
          mt: 2,
          border:
            "1px solid #E7DDD9",
          backgroundColor:
            "#F8F5F3",
          color:
            "#2B211F",
        }}
      >

        A loan is considered at risk
        when it has an unpaid scheduled
        amount past its due date.
        Advanced payments that fully
        cover scheduled dues are not
        included in PAR.

      </Alert>

    </Box>
  );
}


// ========================================
// AGING LABEL
// ========================================

function getAgingLabel(
  daysPastDue: number,
): string {

  if (daysPastDue <= 30) {
    return "1–30 Days";
  }

  if (daysPastDue <= 60) {
    return "31–60 Days";
  }

  if (daysPastDue <= 90) {
    return "61–90 Days";
  }

  return "90+ Days";
}


function getAgingData(
  aging: PARAging[],
  key: PARAging["aging"],
) {
  return (
    aging.find(
      (item) => item.aging === key,
    ) ?? {
      aging: key,
      loans: 0,
      default_amount: 0,
    }
  );
}


// ========================================
// SUMMARY CARD
// ========================================

interface SummaryCardProps {
  title: string;
  value: string;
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
        >

          <Box>

            <Typography
              variant="body2"
              sx={{
                color:
                  "#756B68",
              }}
            >
              {title}
            </Typography>

            <Typography
              variant="h6"
              fontWeight={700}
              sx={{
                mt: 1,
                color:
                  "#2B211F",
              }}
            >
              {value}
            </Typography>

          </Box>


          <Box
            sx={{
              width: 42,
              height: 42,
              display:
                "flex",
              alignItems:
                "center",
              justifyContent:
                "center",
              borderRadius: 2,
              backgroundColor:
                "#F8F5F3",
              color:
                "#8F2115",
            }}
          >

            {icon}

          </Box>

        </Stack>

      </CardContent>

    </Card>
  );
}
