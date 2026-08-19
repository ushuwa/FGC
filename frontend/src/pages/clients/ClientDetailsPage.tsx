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
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  IconButton,
  InputAdornment,
  Stack,
  TextField,
  Typography,
} from "@mui/material";

import {
  Add,
  DeleteOutline,
  EditOutlined,
  Search,
  VisibilityOutlined,
} from "@mui/icons-material";

import {
  DataGrid,
  type GridColDef,
} from "@mui/x-data-grid";

import {
  useForm,
  type SubmitHandler,
} from "react-hook-form";

import {
  useNavigate,
} from "react-router-dom";

import {
  createClient,
  deleteClient,
  getClients,
  updateClient,
} from "../../services/clientService";

import type {
  Client,
  CreateClientRequest,
} from "../../types/client";

interface ClientFormValues {
  first_name: string;
  last_name: string;
  contact_number: string;
  email: string;
  current_address: string;
}

const emptyForm: ClientFormValues = {
  first_name: "",
  last_name: "",
  contact_number: "",
  email: "",
  current_address: "",
};

export default function ClientDetailsPage() {
  const navigate = useNavigate();

  const [clients, setClients] =
    useState<Client[]>([]);

  const [loading, setLoading] =
    useState(true);

  const [error, setError] =
    useState("");

  const [search, setSearch] =
    useState("");

  const [openForm, setOpenForm] =
    useState(false);

  const [editingClient, setEditingClient] =
    useState<Client | null>(null);

  const [deleteTarget, setDeleteTarget] =
    useState<Client | null>(null);

  const [saving, setSaving] =
    useState(false);

  const [deleting, setDeleting] =
    useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: {
      errors,
    },
  } = useForm<ClientFormValues>({
    defaultValues: emptyForm,
  });

  // ============================================
  // LOAD CLIENTS
  // ============================================

  const loadClients = async (
    searchValue = "",
  ) => {
    try {
      setLoading(true);
      setError("");

      const data =
        await getClients(
          searchValue,
        );

      setClients(data);
    } catch (err) {
      console.error(
        "Failed to load clients:",
        err,
      );

      setError(
        "Unable to load clients.",
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadClients("");
  }, []);

  // ============================================
  // SEARCH
  // ============================================

  const handleSearch = () => {
    void loadClients(search);
  };

  // ============================================
  // CREATE
  // ============================================

  const handleOpenCreate = () => {
    setEditingClient(null);

    reset(emptyForm);

    setOpenForm(true);
  };

  // ============================================
  // EDIT
  // ============================================

  const handleOpenEdit = (
    client: Client,
  ) => {
    setEditingClient(client);

    reset({
      first_name:
        client.first_name,

      last_name:
        client.last_name,

      contact_number:
        client.contact_number ??
        "",

      email:
        client.email ??
        "",

      current_address:
        client.current_address ??
        "",
    });

    setOpenForm(true);
  };

  // ============================================
  // CLOSE FORM
  // ============================================

  const handleCloseForm = () => {
    if (saving) {
      return;
    }

    setOpenForm(false);
    setEditingClient(null);

    reset(emptyForm);
  };

  // ============================================
  // SAVE
  // ============================================

  const onSubmit: SubmitHandler<
    ClientFormValues
  > = async (values) => {
    try {
      setSaving(true);
      setError("");

      const payload: CreateClientRequest =
        {
          first_name:
            values.first_name.trim(),

          last_name:
            values.last_name.trim(),

          contact_number:
            values.contact_number.trim(),

          email:
            values.email.trim(),

          current_address:
            values.current_address.trim(),
        };

      if (editingClient) {
        await updateClient(
          editingClient.id,
          payload,
        );
      } else {
        await createClient(
          payload,
        );
      }

      handleCloseForm();

      await loadClients(search);
    } catch (err) {
      console.error(
        "Failed to save client:",
        err,
      );

      setError(
        editingClient
          ? "Unable to update client."
          : "Unable to create client.",
      );
    } finally {
      setSaving(false);
    }
  };

  // ============================================
  // DELETE
  // ============================================

  const handleDelete = async () => {
    if (!deleteTarget) {
      return;
    }

    try {
      setDeleting(true);
      setError("");

      await deleteClient(
        deleteTarget.id,
      );

      setDeleteTarget(null);

      await loadClients(search);
    } catch (err) {
      console.error(
        "Failed to delete client:",
        err,
      );

      setError(
        "Unable to delete client.",
      );
    } finally {
      setDeleting(false);
    }
  };

  // ============================================
  // TABLE COLUMNS
  // ============================================

  const columns: GridColDef[] = [
    {
      field: "id",
      headerName: "ID",
      width: 80,
    },

    {
      field: "name",
      headerName: "Client Name",
      minWidth: 220,
      flex: 1,

      valueGetter: (
        _value,
        row,
      ) =>
        `${row.first_name} ${row.last_name}`,
    },

    {
      field: "contact_number",
      headerName: "Contact Number",
      minWidth: 160,
      flex: 0.8,

      renderCell: (
        params,
      ) =>
        params.value ||
        "—",
    },

    {
      field: "email",
      headerName: "Email",
      minWidth: 200,
      flex: 1,

      renderCell: (
        params,
      ) =>
        params.value ||
        "—",
    },

    {
      field: "current_address",
      headerName: "Address",
      minWidth: 250,
      flex: 1.2,

      renderCell: (
        params,
      ) =>
        params.value ||
        "—",
    },

    {
      field: "actions",
      headerName: "Actions",
      width: 150,
      sortable: false,
      filterable: false,

      renderCell: (
        params,
      ) => (
        <Stack
          direction="row"
          spacing={0.5}
        >
          {/* VIEW */}
          <IconButton
            size="small"
            onClick={() =>
              navigate(
                `/clients/${params.row.id}`,
              )
            }
            sx={{
              color: "#8F2115",
            }}
          >
            <VisibilityOutlined
              fontSize="small"
            />
          </IconButton>

          {/* EDIT */}
          <IconButton
            size="small"
            onClick={() =>
              handleOpenEdit(
                params.row,
              )
            }
            sx={{
              color: "#D0B050",
            }}
          >
            <EditOutlined
              fontSize="small"
            />
          </IconButton>

          {/* DELETE */}
          <IconButton
            size="small"
            onClick={() =>
              setDeleteTarget(
                params.row,
              )
            }
            sx={{
              color: "#A51D1D",
            }}
          >
            <DeleteOutline
              fontSize="small"
            />
          </IconButton>
        </Stack>
      ),
    },
  ];

  // ============================================
  // UI
  // ============================================

  return (
    <Box
      sx={{
        width: "100%",
      }}
    >
      {/* ========================================
          PAGE HEADER
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
            Client Details
          </Typography>

          <Typography
            variant="body2"
            sx={{
              mt: 0.5,
              color: "#756B68",
            }}
          >
            Manage client information
            and records.
          </Typography>
        </Box>

        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={
            handleOpenCreate
          }
          sx={{
            backgroundColor:
              "#8F2115",

            "&:hover": {
              backgroundColor:
                "#70150F",
            },
          }}
        >
          Add Client
        </Button>
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
          CLIENT TABLE CARD
      ======================================== */}

      <Card>
        <CardContent>
          {/* SEARCH */}

          <Stack
            direction={{
              xs: "column",
              sm: "row",
            }}
            spacing={1.5}
            sx={{
              mb: 2.5,
            }}
          >
            <TextField
              fullWidth
              value={search}
              placeholder="Search by name, contact, or email"
              onChange={(event) =>
                setSearch(
                  event.target.value,
                )
              }
              onKeyDown={(event) => {
                if (
                  event.key ===
                  "Enter"
                ) {
                  handleSearch();
                }
              }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search
                      sx={{
                        color:
                          "#8F2115",
                      }}
                    />
                  </InputAdornment>
                ),
              }}
            />

            <Button
              variant="outlined"
              onClick={
                handleSearch
              }
              sx={{
                minWidth: 110,

                borderColor:
                  "#8F2115",

                color:
                  "#8F2115",

                "&:hover": {
                  borderColor:
                    "#70150F",

                  backgroundColor:
                    "rgba(143,33,21,0.05)",
                },
              }}
            >
              Search
            </Button>
          </Stack>

          {/* DATA GRID */}

          <Box
            sx={{
              width: "100%",
            }}
          >
            <DataGrid
              rows={clients}
              columns={columns}
              loading={loading}
              autoHeight
              disableRowSelectionOnClick
              pageSizeOptions={[
                10,
                25,
                50,
              ]}
              initialState={{
                pagination: {
                  paginationModel: {
                    pageSize: 10,
                    page: 0,
                  },
                },
              }}
              sx={{
                border: "none",

                "& .MuiDataGrid-columnHeaders":
                  {
                    backgroundColor:
                      "#F8F5F3",

                    borderBottom:
                      "1px solid #E7DDD9",
                  },

                "& .MuiDataGrid-columnHeaderTitle":
                  {
                    fontWeight: 700,

                    color:
                      "#2B211F",
                  },

                "& .MuiDataGrid-row:hover":
                  {
                    backgroundColor:
                      "rgba(208,176,80,0.06)",
                  },

                "& .MuiDataGrid-cell":
                  {
                    borderBottom:
                      "1px solid #F0EAE7",
                  },
              }}
            />
          </Box>
        </CardContent>
      </Card>

      {/* ========================================
          CREATE / EDIT DIALOG
      ======================================== */}

      <Dialog
        open={openForm}
        onClose={
          handleCloseForm
        }
        fullWidth
        maxWidth="sm"
      >
        <form
          onSubmit={handleSubmit(
            onSubmit,
          )}
        >
          <DialogTitle
            sx={{
              fontWeight: 700,
              color: "#2B211F",
            }}
          >
            {editingClient
              ? "Edit Client"
              : "Add Client"}
          </DialogTitle>

          <DialogContent>
            <Stack
              spacing={2}
              sx={{
                pt: 1,
              }}
            >
              <Grid
                container
                spacing={2}
              >
                <Grid
                  size={{
                  xs: 12,
                  sm: 6
                }}
                  
                >
                  <TextField
                    fullWidth
                    label="First Name"
                    {...register(
                      "first_name",
                      {
                        required:
                          "First name is required",
                      },
                    )}
                    error={
                      !!errors.first_name
                    }
                    helperText={
                      errors.first_name
                        ?.message
                    }
                  />
                </Grid>

                <Grid
                  size={{
                  xs: 12,
                  sm: 6
                }}
                >

                  <TextField
                    fullWidth
                    label="Last Name"
                    {...register(
                      "last_name",
                      {
                        required:
                          "Last name is required",
                      },
                    )}
                    error={
                      !!errors.last_name
                    }
                    helperText={
                      errors.last_name
                        ?.message
                    }
                  />
                </Grid>
              </Grid>

              <TextField
                fullWidth
                label="Contact Number"
                {...register(
                  "contact_number",
                )}
              />

              <TextField
                fullWidth
                label="Email"
                type="email"
                {...register(
                  "email",
                )}
              />

              <TextField
                fullWidth
                label="Current Address"
                multiline
                rows={3}
                {...register(
                  "current_address",
                )}
              />
            </Stack>
          </DialogContent>

          <DialogActions
            sx={{
              px: 3,
              pb: 2.5,
            }}
          >
            <Button
              onClick={
                handleCloseForm
              }
              disabled={saving}
              sx={{
                color: "#756B68",
              }}
            >
              Cancel
            </Button>

            <Button
              type="submit"
              variant="contained"
              disabled={saving}
              sx={{
                minWidth: 110,

                backgroundColor:
                  "#8F2115",

                "&:hover": {
                  backgroundColor:
                    "#70150F",
                },
              }}
            >
              {saving
                ? "Saving..."
                : editingClient
                  ? "Update"
                  : "Save"}
            </Button>
          </DialogActions>
        </form>
      </Dialog>

      {/* ========================================
          DELETE CONFIRMATION
      ======================================== */}

      <Dialog
        open={
          deleteTarget !== null
        }
        onClose={() => {
          if (!deleting) {
            setDeleteTarget(null);
          }
        }}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle
          sx={{
            fontWeight: 700,
          }}
        >
          Delete Client?
        </DialogTitle>

        <DialogContent>
          <Typography>
            Are you sure you want
            to delete{" "}
            <strong>
              {
                deleteTarget?.first_name
              }{" "}
              {
                deleteTarget?.last_name
              }
            </strong>
            ?
          </Typography>

          <Typography
            variant="body2"
            sx={{
              mt: 1,
              color: "#756B68",
            }}
          >
            This action cannot be
            undone.
          </Typography>
        </DialogContent>

        <DialogActions
          sx={{
            px: 3,
            pb: 2.5,
          }}
        >
          <Button
            onClick={() =>
              setDeleteTarget(null)
            }
            disabled={deleting}
          >
            Cancel
          </Button>

          <Button
            variant="contained"
            onClick={
              handleDelete
            }
            disabled={deleting}
            sx={{
              backgroundColor:
                "#A51D1D",

              "&:hover": {
                backgroundColor:
                  "#851515",
              },
            }}
          >
            {deleting
              ? "Deleting..."
              : "Delete"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}