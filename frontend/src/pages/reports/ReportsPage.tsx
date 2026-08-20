import {
  useEffect,
  useState,
} from "react";

import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  MenuItem,
  Select,
  Stack,
  TextField,
  Typography,
} from "@mui/material";

import {
  PictureAsPdfOutlined,
  PrintOutlined,
} from "@mui/icons-material";

import {
  downloadCollectionReport,
  downloadAmortizationReport,
  downloadPARReport,
  downloadLoanPortfolioReport,
  downloadLoanMaturityReport,
} from "../../services/reportService";

import {
  getLoans,
} from "../../services/loanService";


// ========================================
// TYPES
// ========================================

interface LoanOption {
  id: number;
  pn_number: string;
  client_name: string;
}


// ========================================
// PAGE
// ========================================

export default function ReportsPage() {

  // ========================================
  // COLLECTION REPORT
  // ========================================

  const [dateFrom, setDateFrom] =
    useState("");

  const [dateTo, setDateTo] =
    useState("");

  const [loading, setLoading] =
    useState(false);

  const [parLoading, setParLoading] =
  useState(false);

  const [portfolioLoading, setPortfolioLoading] =
  useState(false);

  const [maturityLoading, setMaturityLoading] =
  useState(false);

  // ========================================
  // AMORTIZATION REPORT
  // ========================================

  const [selectedLoanId, setSelectedLoanId] =
    useState<number | "">("");

  const [loans, setLoans] =
    useState<LoanOption[]>([]);

  const [loansLoading, setLoansLoading] =
    useState(true);

  const [amortizationLoading, setAmortizationLoading] =
    useState(false);


  // ========================================
  // ERROR
  // ========================================

  const [error, setError] =
    useState("");


  // ========================================
  // LOAD LOANS
  // ========================================

  useEffect(() => {

    const loadLoans = async () => {

      try {

        setLoansLoading(true);

        setError("");

        const data =
          await getLoans();

        setLoans(
          data ?? [],
        );

      } catch (err) {

        console.error(
          "Failed to load loans:",
          err,
        );

        setLoans([]);

        setError(
          "Unable to load loans.",
        );

      } finally {

        setLoansLoading(false);

      }
    };

    void loadLoans();

  }, []);


  // ========================================
  // COLLECTION PDF
  // ========================================

  const generatePDF =
    async () => {

      if (!dateFrom || !dateTo) {

        setError(
          "Please select a date range.",
        );

        return;
      }

      if (dateFrom > dateTo) {

        setError(
          "Date From cannot be later than Date To.",
        );

        return;
      }

      try {

        setLoading(true);

        setError("");

        const blob =
          await downloadCollectionReport(
            dateFrom,
            dateTo,
          );

        const url =
          window.URL.createObjectURL(
            blob,
          );

        const link =
          document.createElement("a");

        link.href = url;

        link.download =
          `summary_collection_${dateFrom}_${dateTo}.pdf`;

        document.body.appendChild(link);

        link.click();

        link.remove();

        window.URL.revokeObjectURL(
          url,
        );

      } catch (err) {

        console.error(
          "Failed to generate collection report:",
          err,
        );

        setError(
          "Unable to generate report.",
        );

      } finally {

        setLoading(false);

      }
    };


  // ========================================
  // PRINT COLLECTION PDF
  // ========================================

  const printPDF =
    async () => {

      if (!dateFrom || !dateTo) {

        setError(
          "Please select a date range.",
        );

        return;
      }

      if (dateFrom > dateTo) {

        setError(
          "Date From cannot be later than Date To.",
        );

        return;
      }

      try {

        setLoading(true);

        setError("");

        const blob =
          await downloadCollectionReport(
            dateFrom,
            dateTo,
          );

        const url =
          window.URL.createObjectURL(
            blob,
          );

        const printWindow =
          window.open(
            url,
            "_blank",
          );

        if (!printWindow) {

          window.URL.revokeObjectURL(
            url,
          );

          setError(
            "Unable to open print preview.",
          );

          return;
        }

        printWindow.onload = () => {

          printWindow.focus();

          printWindow.print();

        };

      } catch (err) {

        console.error(
          "Failed to print collection report:",
          err,
        );

        setError(
          "Unable to print report.",
        );

      } finally {

        setLoading(false);

      }
    };


  // ========================================
  // AMORTIZATION PDF
  // ========================================

  const generateAmortizationPDF =
    async () => {

      if (!selectedLoanId) {

        setError(
          "Please select a loan.",
        );

        return;
      }

      try {

        setAmortizationLoading(true);

        setError("");

        const blob =
          await downloadAmortizationReport(
            selectedLoanId,
          );

        const url =
          window.URL.createObjectURL(
            blob,
          );

        const link =
          document.createElement("a");

        link.href = url;

        link.download =
          "loan_amortization.pdf";

        document.body.appendChild(link);

        link.click();

        link.remove();

        window.URL.revokeObjectURL(
          url,
        );

      } catch (err) {

        console.error(
          "Failed to generate amortization report:",
          err,
        );

        setError(
          "Unable to generate loan amortization report.",
        );

      } finally {

        setAmortizationLoading(false);

      }
    };



    //PAR
    const generatePARPDF =
  async () => {

    try {

      setParLoading(true);

      setError("");

      const blob =
        await downloadPARReport();

      const url =
        window.URL.createObjectURL(
          blob,
        );

      const link =
        document.createElement("a");

      link.href = url;

      link.download =
        "portfolio_at_risk.pdf";

      document.body.appendChild(link);

      link.click();

      link.remove();

      window.URL.revokeObjectURL(
        url,
      );

    } catch (err) {

      console.error(
        "Failed to generate PAR report:",
        err,
      );

      setError(
        "Unable to generate Portfolio at Risk report.",
      );

    } finally {

      setParLoading(false);

    }
  };


  //PORTFOLIO
  const generateLoanPortfolioPDF =
  async () => {

    try {

      setPortfolioLoading(true);

      setError("");

      const blob =
        await downloadLoanPortfolioReport();

      const url =
        window.URL.createObjectURL(
          blob,
        );

      const link =
        document.createElement("a");

      link.href = url;

      link.download =
        "loan_portfolio_summary.pdf";

      document.body.appendChild(link);

      link.click();

      link.remove();

      window.URL.revokeObjectURL(
        url,
      );

    } catch (err) {

      console.error(
        "Failed to generate loan portfolio report:",
        err,
      );

      setError(
        "Unable to generate Loan Portfolio Summary report.",
      );

    } finally {

      setPortfolioLoading(false);

    }
  };


  //MATURITY
  const generateLoanMaturityPDF =
  async () => {

    try {

      setMaturityLoading(true);

      setError("");

      const blob =
        await downloadLoanMaturityReport();

      const url =
        window.URL.createObjectURL(
          blob,
        );

      const link =
        document.createElement("a");

      link.href = url;

      link.download =
        "loan_maturity_due.pdf";

      document.body.appendChild(link);

      link.click();

      link.remove();

      window.URL.revokeObjectURL(
        url,
      );

    } catch (err) {

      console.error(
        "Failed to generate loan maturity report:",
        err,
      );

      setError(
        "Unable to generate Loan Maturity / Due report.",
      );

    } finally {

      setMaturityLoading(false);

    }
  };



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
          Reports
        </Typography>

        <Typography
          variant="body2"
          sx={{
            color: "#756B68",
          }}
        >
          Generate official FGC reports.
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
          onClose={() =>
            setError("")
          }
        >
          {error}
        </Alert>
      )}


      {/* ========================================
          SUMMARY OF COLLECTION
      ======================================== */}

      <Card>

        <CardContent>

          <Typography
            variant="h6"
            fontWeight={700}
            sx={{
              color: "#2B211F",
              mb: 3,
            }}
          >
            Summary of Collection
          </Typography>

          <Stack
            direction={{
              xs: "column",
              md: "row",
            }}
            spacing={2}
          >

            <TextField
              type="date"
              label="Date From"
              value={dateFrom}
              onChange={(event) => {

                setDateFrom(
                  event.target.value,
                );

                setError("");

              }}
              InputLabelProps={{
                shrink: true,
              }}
            />


            <TextField
              type="date"
              label="Date To"
              value={dateTo}
              onChange={(event) => {

                setDateTo(
                  event.target.value,
                );

                setError("");

              }}
              InputLabelProps={{
                shrink: true,
              }}
            />


            <Button
              variant="contained"
              startIcon={
                <PictureAsPdfOutlined />
              }
              disabled={
                loading ||
                !dateFrom ||
                !dateTo
              }
              onClick={
                generatePDF
              }
              sx={{
                backgroundColor:
                  "#8F2115",

                "&:hover": {
                  backgroundColor:
                    "#741B12",
                },
              }}
            >
              {loading
                ? "Generating..."
                : "Generate PDF"}
            </Button>


            <Button
              variant="outlined"
              startIcon={
                <PrintOutlined />
              }
              disabled={
                loading ||
                !dateFrom ||
                !dateTo
              }
              onClick={
                printPDF
              }
              sx={{
                color:
                  "#8F2115",

                borderColor:
                  "#8F2115",

                "&:hover": {
                  borderColor:
                    "#741B12",

                  backgroundColor:
                    "#F8F5F3",
                },
              }}
            >
              {loading
                ? "Preparing..."
                : "Print"}
            </Button>

          </Stack>

        </CardContent>

      </Card>


      {/* ========================================
          LOAN ACCOUNT AMORTIZATION
      ======================================== */}

      <Card
        sx={{
          mt: 2,
        }}
      >

        <CardContent>

          <Typography
            variant="h6"
            fontWeight={700}
            sx={{
              color: "#2B211F",
              mb: 1,
            }}
          >
            Loan Account Amortization
          </Typography>


          <Typography
            variant="body2"
            sx={{
              color: "#756B68",
              mb: 3,
            }}
          >
            Generate the complete amortization
            schedule and payment status for a
            specific loan.
          </Typography>


          <Stack
            direction={{
              xs: "column",
              md: "row",
            }}
            spacing={2}
            alignItems={{
              xs: "stretch",
              md: "center",
            }}
          >

            {/* ========================================
                LOAN SELECT
            ======================================== */}

            <Select
              value={selectedLoanId}
              onChange={(event) => {

                const value =
                  event.target.value;

                setSelectedLoanId(
                  value === null
                    ? ""
                    : Number(value),
                );

                setError("");

              }}
              displayEmpty
              disabled={
                loansLoading ||
                loans.length === 0
              }
              sx={{
                minWidth: {
                  xs: "100%",
                  md: 360,
                },
              }}
            >

              <MenuItem value="">
                {loansLoading
                  ? "Loading loans..."
                  : loans.length === 0
                    ? "No loans available"
                    : "Select Loan / PN"}
              </MenuItem>


              {loans.map(
                (loan) => (

                  <MenuItem
                    key={loan.id}
                    value={loan.id}
                  >
                    {loan.pn_number}
                    {" - "}
                    {loan.client_name}
                  </MenuItem>

                ),
              )}

            </Select>


            {/* ========================================
                GENERATE
            ======================================== */}

            <Button
              variant="contained"
              startIcon={
                <PictureAsPdfOutlined />
              }
              disabled={
                selectedLoanId === "" ||
                amortizationLoading ||
                loansLoading
              }
              onClick={
                generateAmortizationPDF
              }
              sx={{
                backgroundColor:
                  "#8F2115",

                "&:hover": {
                  backgroundColor:
                    "#741B12",
                },
              }}
            >
              {amortizationLoading
                ? "Generating..."
                : "Generate PDF"}
            </Button>

          </Stack>

        </CardContent>

      </Card>


       {/* ========================================
          PAR
      ======================================== */}

      <Card
  sx={{
    mt: 2,
  }}
>
  <CardContent>

    <Typography
      variant="h6"
      fontWeight={700}
      sx={{
        color: "#2B211F",
        mb: 1,
      }}
    >
      Portfolio at Risk
    </Typography>

    <Typography
      variant="body2"
      sx={{
        color: "#756B68",
        mb: 3,
      }}
    >
      Generate a report of loans with unpaid
      scheduled amounts past their due date.
    </Typography>

    <Button
      variant="contained"
      startIcon={
        <PictureAsPdfOutlined />
      }
      disabled={parLoading}
      onClick={generatePARPDF}
      sx={{
        backgroundColor:
          "#8F2115",

        "&:hover": {
          backgroundColor:
            "#741B12",
        },
      }}
    >
      {parLoading
        ? "Generating..."
        : "Generate PDF"}
    </Button>

  </CardContent>
</Card>

   {/* ========================================
          PORFOLIO
      ======================================== */}

<Card
  sx={{
    mt: 2,
  }}
>
  <CardContent>

    <Typography
      variant="h6"
      fontWeight={700}
      sx={{
        color: "#2B211F",
        mb: 1,
      }}
    >
      Loan Portfolio Summary
    </Typography>

    <Typography
      variant="body2"
      sx={{
        color: "#756B68",
        mb: 3,
      }}
    >
      Generate a complete summary of the
      current loan portfolio, including
      outstanding balances, loan status,
      and loan type breakdown.
    </Typography>

    <Button
      variant="contained"
      startIcon={
        <PictureAsPdfOutlined />
      }
      disabled={portfolioLoading}
      onClick={
        generateLoanPortfolioPDF
      }
      sx={{
        backgroundColor:
          "#8F2115",

        "&:hover": {
          backgroundColor:
            "#741B12",
        },
      }}
    >
      {portfolioLoading
        ? "Generating..."
        : "Generate PDF"}
    </Button>

  </CardContent>
</Card>



<Card
  sx={{
    mt: 2,
  }}
>
  <CardContent>

    <Typography
      variant="h6"
      fontWeight={700}
      sx={{
        color: "#2B211F",
        mb: 1,
      }}
    >
      Loan Maturity / Due Report
    </Typography>

    <Typography
      variant="body2"
      sx={{
        color: "#756B68",
        mb: 3,
      }}
    >
      Generate a report of loan maturity dates,
      upcoming maturities, matured loans, and
      outstanding balances.
    </Typography>

    <Button
      variant="contained"
      startIcon={
        <PictureAsPdfOutlined />
      }
      disabled={maturityLoading}
      onClick={
        generateLoanMaturityPDF
      }
      sx={{
        backgroundColor:
          "#8F2115",

        "&:hover": {
          backgroundColor:
            "#741B12",
        },
      }}
    >
      {maturityLoading
        ? "Generating..."
        : "Generate PDF"}
    </Button>

  </CardContent>
</Card>






    </Box>
  );
}