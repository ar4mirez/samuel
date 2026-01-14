# R Guide

> **Applies to**: R 4.0+, RStudio, Tidyverse, Shiny, Statistical Computing, Data Science

---

## Core Principles

1. **Reproducibility**: Scripts, not console commands
2. **Tidyverse First**: Modern, consistent R programming
3. **Vectorization**: Avoid explicit loops when possible
4. **Functional Style**: Pure functions, immutable data
5. **Documentation**: Roxygen2 for packages, comments for scripts

---

## Language-Specific Guardrails

### R Version & Setup
- ✓ Use R 4.0+ (4.3+ recommended)
- ✓ Use `renv` for dependency management
- ✓ Use RStudio or VSCode with R extension
- ✓ Pin package versions in `renv.lock`

### Code Style (tidyverse style guide)
- ✓ Use snake_case for variables and functions
- ✓ Use SCREAMING_SNAKE_CASE for constants
- ✓ 2-space indentation
- ✓ Max line length: 80 characters
- ✓ Use `<-` for assignment (not `=`)
- ✓ Space after commas, around operators
- ✓ Use `styler` package for formatting
- ✓ Use `lintr` for linting

### Data Handling
- ✓ Use tibbles over data.frames
- ✓ Use tidyverse verbs for data manipulation
- ✓ Use the pipe operator `|>` (R 4.1+) or `%>%`
- ✓ Prefer explicit column names over indices
- ✓ Handle missing values explicitly

### Functions
- ✓ Functions should do one thing well
- ✓ Use meaningful parameter names
- ✓ Document with Roxygen2 comments
- ✓ Return values explicitly
- ✓ Avoid side effects

### Packages
- ✓ Document packages with Roxygen2
- ✓ Use testthat for testing
- ✓ Follow CRAN submission guidelines
- ✓ Include NEWS.md for changelog
- ✓ Use pkgdown for documentation websites

---

## Project Structure

### Analysis Project
```
myproject/
├── R/                    # R scripts/functions
│   ├── 01_load_data.R
│   ├── 02_clean_data.R
│   ├── 03_analysis.R
│   └── functions/        # Helper functions
├── data/                 # Raw data (read-only)
│   └── raw/
├── data-raw/            # Data processing scripts
├── output/              # Generated outputs
│   ├── figures/
│   └── tables/
├── reports/             # R Markdown reports
│   └── analysis.Rmd
├── tests/               # Tests
│   └── testthat/
├── renv/                # renv library
├── renv.lock            # Package versions
├── .Rprofile           # Project settings
├── _targets.R          # targets pipeline (optional)
└── README.md
```

### Package Structure
```
mypackage/
├── R/                   # R source files
│   └── functions.R
├── man/                 # Documentation (generated)
├── tests/
│   └── testthat/
├── vignettes/          # Long-form documentation
├── data/               # Package data
├── inst/               # Additional files
├── DESCRIPTION         # Package metadata
├── NAMESPACE           # Exports (generated)
├── NEWS.md             # Changelog
└── README.md
```

---

## Data Manipulation (tidyverse)

### Loading Data
```r
library(tidyverse)

# CSV
data <- read_csv("data/file.csv")

# Excel
library(readxl)
data <- read_excel("data/file.xlsx", sheet = "Sheet1")

# Database
library(DBI)
con <- dbConnect(RSQLite::SQLite(), "database.db")
data <- dbGetQuery(con, "SELECT * FROM table_name")
dbDisconnect(con)

# Parquet (for large data)
library(arrow)
data <- read_parquet("data/file.parquet")
```

### Core dplyr Verbs
```r
library(dplyr)

# Select columns
data |>
  select(name, age, email)

# Filter rows
data |>
  filter(age >= 18, status == "active")

# Create/modify columns
data |>
  mutate(
    age_group = case_when(
      age < 18 ~ "minor",
      age < 65 ~ "adult",
      TRUE ~ "senior"
    ),
    full_name = paste(first_name, last_name)
  )

# Summarize
data |>
  group_by(department) |>
  summarize(
    count = n(),
    avg_salary = mean(salary, na.rm = TRUE),
    total_salary = sum(salary, na.rm = TRUE)
  )

# Sort
data |>
  arrange(desc(created_at), name)

# Join
users |>
  left_join(orders, by = "user_id") |>
  inner_join(products, by = "product_id")

# Pivot
# Long to wide
data |>
  pivot_wider(
    names_from = year,
    values_from = value
  )

# Wide to long
data |>
  pivot_longer(
    cols = starts_with("year_"),
    names_to = "year",
    values_to = "value"
  )
```

### Data Cleaning
```r
library(tidyr)
library(stringr)

# Handle missing values
data |>
  drop_na(important_column) |>          # Remove rows with NA
  replace_na(list(value = 0)) |>        # Replace NA with value
  fill(category, .direction = "down")   # Fill NA with previous value

# String manipulation
data |>
  mutate(
    email_lower = str_to_lower(email),
    name_clean = str_trim(str_squish(name)),
    domain = str_extract(email, "@(.+)$")
  ) |>
  filter(str_detect(email, "@company\\.com$"))

# Type conversion
data |>
  mutate(
    date = as.Date(date_string, format = "%Y-%m-%d"),
    amount = as.numeric(amount_string),
    category = as.factor(category)
  )
```

---

## Functions

### Basic Function
```r
#' Calculate the mean of positive values
#'
#' @param x A numeric vector
#' @param na.rm Logical. Should NA values be removed?
#' @return The mean of positive values in x
#' @examples
#' mean_positive(c(-1, 2, 3, NA, 5))
#' @export
mean_positive <- function(x, na.rm = TRUE) {
  if (!is.numeric(x)) {
    stop("x must be numeric")
  }

  positive_values <- x[x > 0]

  if (length(positive_values) == 0) {
    return(NA_real_)
  }

  mean(positive_values, na.rm = na.rm)
}
```

### Defensive Programming
```r
#' Process user data
#'
#' @param data A data frame with user information
#' @param min_age Minimum age filter (default: 0)
#' @return Processed data frame
process_users <- function(data, min_age = 0) {
  # Input validation
  stopifnot(
    is.data.frame(data),
    is.numeric(min_age),
    min_age >= 0
  )

  required_cols <- c("user_id", "name", "age")
  missing_cols <- setdiff(required_cols, names(data))

  if (length(missing_cols) > 0) {
    stop(
      "Missing required columns: ",
      paste(missing_cols, collapse = ", ")
    )
  }

  # Process data
  data |>
    filter(age >= min_age) |>
    mutate(
      name = str_to_title(name),
      age_group = cut(
        age,
        breaks = c(0, 18, 65, Inf),
        labels = c("minor", "adult", "senior")
      )
    )
}
```

### Functional Programming with purrr
```r
library(purrr)

# Map over list
files <- list.files("data/", pattern = "\\.csv$", full.names = TRUE)
all_data <- map(files, read_csv) |>
  list_rbind()

# Map with type specification
numbers <- list(1:3, 4:6, 7:9)
sums <- map_dbl(numbers, sum)

# Map over multiple inputs
map2(names, ages, ~ paste(.x, "is", .y, "years old"))

# Conditional operations
data |>
  mutate(
    processed = map_if(
      values,
      is.numeric,
      ~ .x * 2
    )
  )

# Safe operations (handle errors)
safe_read <- safely(read_csv)
results <- map(files, safe_read)
successful <- map(results, "result") |> compact()
errors <- map(results, "error") |> compact()

# Reduce
list(df1, df2, df3) |>
  reduce(left_join, by = "id")
```

---

## Data Visualization (ggplot2)

### Basic Plots
```r
library(ggplot2)

# Scatter plot
ggplot(data, aes(x = age, y = income)) +
  geom_point(alpha = 0.5) +
  geom_smooth(method = "lm") +
  labs(
    title = "Income vs Age",
    x = "Age (years)",
    y = "Income ($)"
  ) +
  theme_minimal()

# Bar plot
ggplot(data, aes(x = category, fill = status)) +
  geom_bar(position = "dodge") +
  scale_fill_brewer(palette = "Set2") +
  coord_flip() +
  theme_minimal()

# Histogram
ggplot(data, aes(x = value)) +
  geom_histogram(bins = 30, fill = "steelblue", color = "white") +
  facet_wrap(~ group) +
  theme_minimal()

# Line plot
ggplot(data, aes(x = date, y = value, color = category)) +
  geom_line(linewidth = 1) +
  scale_x_date(date_breaks = "1 month", date_labels = "%b %Y") +
  theme_minimal() +
  theme(axis.text.x = element_text(angle = 45, hjust = 1))

# Save plot
ggsave(
  "output/figures/my_plot.png",
  width = 10,
  height = 6,
  dpi = 300
)
```

### Custom Theme
```r
theme_custom <- function() {
  theme_minimal() +
    theme(
      text = element_text(family = "Arial"),
      plot.title = element_text(size = 14, face = "bold"),
      axis.title = element_text(size = 12),
      legend.position = "bottom",
      panel.grid.minor = element_blank()
    )
}

# Use custom theme
ggplot(data, aes(x, y)) +
  geom_point() +
  theme_custom()
```

---

## Statistical Analysis

### Descriptive Statistics
```r
# Summary
summary(data$value)

# Grouped summary
data |>
  group_by(category) |>
  summarize(
    n = n(),
    mean = mean(value, na.rm = TRUE),
    sd = sd(value, na.rm = TRUE),
    median = median(value, na.rm = TRUE),
    q25 = quantile(value, 0.25, na.rm = TRUE),
    q75 = quantile(value, 0.75, na.rm = TRUE)
  )

# Correlation matrix
cor_matrix <- data |>
  select(where(is.numeric)) |>
  cor(use = "complete.obs")
```

### Statistical Tests
```r
# t-test
t.test(value ~ group, data = data)

# Chi-square test
chisq.test(table(data$category1, data$category2))

# ANOVA
aov_result <- aov(value ~ group1 * group2, data = data)
summary(aov_result)

# Linear regression
model <- lm(outcome ~ predictor1 + predictor2, data = data)
summary(model)

# Using broom for tidy output
library(broom)
tidy(model)      # Coefficients
glance(model)    # Model statistics
augment(model)   # Predictions and residuals
```

### Machine Learning with tidymodels
```r
library(tidymodels)

# Split data
set.seed(123)
data_split <- initial_split(data, prop = 0.8, strata = outcome)
train_data <- training(data_split)
test_data <- testing(data_split)

# Create recipe (preprocessing)
recipe <- recipe(outcome ~ ., data = train_data) |>
  step_normalize(all_numeric_predictors()) |>
  step_dummy(all_nominal_predictors())

# Define model
model_spec <- logistic_reg() |>
  set_engine("glm") |>
  set_mode("classification")

# Create workflow
workflow <- workflow() |>
  add_recipe(recipe) |>
  add_model(model_spec)

# Fit model
fitted_model <- fit(workflow, data = train_data)

# Predictions
predictions <- predict(fitted_model, test_data) |>
  bind_cols(test_data)

# Evaluate
metrics(predictions, truth = outcome, estimate = .pred_class)
```

---

## Testing with testthat

### Test Structure
```r
# tests/testthat/test-functions.R
library(testthat)

test_that("mean_positive calculates correctly", {
  expect_equal(mean_positive(c(1, 2, 3)), 2)
  expect_equal(mean_positive(c(-1, 2, 3)), 2.5)
  expect_equal(mean_positive(c(-1, -2, -3)), NA_real_)
})

test_that("mean_positive handles NA values", {
  expect_equal(mean_positive(c(1, 2, NA)), 1.5)
  expect_equal(mean_positive(c(NA, NA)), NA_real_)
})

test_that("mean_positive validates input", {
  expect_error(mean_positive("not numeric"))
  expect_error(mean_positive(NULL))
})

test_that("process_users filters by age", {
  test_data <- tibble(
    user_id = 1:3,
    name = c("alice", "bob", "charlie"),
    age = c(15, 25, 70)
  )

  result <- process_users(test_data, min_age = 18)

  expect_equal(nrow(result), 2)
  expect_equal(result$name, c("Bob", "Charlie"))
})
```

### Running Tests
```r
# Run all tests
devtools::test()

# Run specific test file
testthat::test_file("tests/testthat/test-functions.R")

# Check package
devtools::check()
```

---

## R Markdown Reports

### Basic R Markdown
```markdown
---
title: "Analysis Report"
author: "Your Name"
date: "`r Sys.Date()`"
output:
  html_document:
    toc: true
    toc_float: true
    code_folding: hide
---

```{r setup, include=FALSE}
knitr::opts_chunk$set(
  echo = TRUE,
  message = FALSE,
  warning = FALSE,
  fig.width = 10,
  fig.height = 6
)

library(tidyverse)
```

## Introduction

This report analyzes...

```{r load-data}
data <- read_csv("data/analysis_data.csv")
```

## Summary Statistics

```{r summary}
data |>
  group_by(category) |>
  summarize(
    n = n(),
    mean = mean(value)
  ) |>
  knitr::kable()
```

## Visualization

```{r plot, fig.cap="Distribution of values"}
ggplot(data, aes(x = value, fill = category)) +
  geom_histogram(bins = 30) +
  facet_wrap(~ category) +
  theme_minimal()
```
```

---

## Shiny Applications

### Basic Shiny App
```r
library(shiny)
library(tidyverse)

ui <- fluidPage(
  titlePanel("Data Explorer"),

  sidebarLayout(
    sidebarPanel(
      selectInput(
        "variable",
        "Select Variable:",
        choices = NULL
      ),
      sliderInput(
        "bins",
        "Number of Bins:",
        min = 10,
        max = 100,
        value = 30
      )
    ),

    mainPanel(
      plotOutput("histogram"),
      tableOutput("summary")
    )
  )
)

server <- function(input, output, session) {
  # Reactive data
  data <- reactive({
    read_csv("data/data.csv")
  })

  # Update choices based on data
  observe({
    numeric_cols <- data() |>
      select(where(is.numeric)) |>
      names()

    updateSelectInput(session, "variable", choices = numeric_cols)
  })

  # Histogram
  output$histogram <- renderPlot({
    req(input$variable)

    ggplot(data(), aes(x = .data[[input$variable]])) +
      geom_histogram(bins = input$bins, fill = "steelblue") +
      theme_minimal()
  })

  # Summary table
  output$summary <- renderTable({
    req(input$variable)

    data() |>
      summarize(
        Mean = mean(.data[[input$variable]], na.rm = TRUE),
        SD = sd(.data[[input$variable]], na.rm = TRUE),
        Min = min(.data[[input$variable]], na.rm = TRUE),
        Max = max(.data[[input$variable]], na.rm = TRUE)
      )
  })
}

shinyApp(ui, server)
```

---

## Performance Optimization

### Tips
```r
# Use data.table for large datasets
library(data.table)
dt <- fread("large_file.csv")
dt[category == "A", .(mean_value = mean(value)), by = group]

# Use vectorized operations
# Bad
result <- c()
for (i in 1:nrow(data)) {
  result <- c(result, data$value[i] * 2)
}

# Good
result <- data$value * 2

# Parallel processing
library(furrr)
plan(multisession, workers = 4)
results <- future_map(items, process_item)

# Profile code
profvis::profvis({
  # Code to profile
})
```

---

## Common Pitfalls

### Avoid These
```r
# T and F instead of TRUE and FALSE
# (T and F can be overwritten)
result <- if (condition) T else F  # Bad

# Floating point comparison
0.1 + 0.2 == 0.3  # FALSE!

# Forgetting drop = FALSE
df[1, ]  # Returns vector if single column

# Not specifying stringsAsFactors
# (default changed in R 4.0)
```

### Do This Instead
```r
# Use full TRUE/FALSE
result <- if (condition) TRUE else FALSE

# Use near() for floating point
near(0.1 + 0.2, 0.3)  # TRUE

# Preserve data frame structure
df[1, , drop = FALSE]

# Be explicit about factors
read_csv("data.csv")  # Keeps as character by default
```

---

## References

- [R for Data Science (2e)](https://r4ds.hadley.nz/)
- [Advanced R](https://adv-r.hadley.nz/)
- [R Packages](https://r-pkgs.org/)
- [tidyverse Style Guide](https://style.tidyverse.org/)
- [Mastering Shiny](https://mastering-shiny.org/)
- [R Graphics Cookbook](https://r-graphics.org/)
- [CRAN](https://cran.r-project.org/)
