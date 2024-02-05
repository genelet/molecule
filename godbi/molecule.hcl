  atoms "m_a" {
    tableName = "m_a"
    pks       = ["id"]
    idAuto    = "id"
    uniques   = ["x", "y"]
    columns "x" {
      typeName    = "string"
      columnLabel = "x"
      notnull     = true
    }
    columns "y" {
      typeName    = "string"
      columnLabel = "y"
      notnull     = true
    }
    columns "z" {
      typeName    = "string"
      columnLabel = "z"
    }
    columns "id" {
      typeName    = "int"
      columnLabel = "id"
      auto        = true
    }
    actions "insert" {
      nextpages "m_b" "insert" {
        relateArgs = {
          id = "id"
        }
      }
    }
    actions "update" {

    }
    actions "insupd" {
      nextpages "m_b" "insert" {
        relateArgs = {
          id = "id"
        }
      }
    }
    actions "delete" {

    }
    actions "delecs" {

    }
    actions "topics" {
      nextpages "m_a" "edit" {
        relateExtra = {
          id = "id"
        }
      }
    }
    actions "edit" {
      nextpages "m_b" "topics" {
        relateExtra = {
          id = "id"
        }
      }
    }
    actions "stmt" {

    }
  }
  atoms "m_b" {
    tableName = "m_b"
    pks       = ["tid"]
    idAuto    = "tid"
    columns "tid" {
      typeName    = "int"
      columnLabel = "tid"
      notnull     = true
      auto        = true
    }
    columns "child" {
      typeName    = "string"
      columnLabel = "child"
    }
    columns "id" {
      typeName    = "int"
      columnLabel = "id"
      notnull     = true
    }
    actions "insert" {

    }
    actions "update" {

    }
    actions "insupd" {

    }
    actions "delete" {

    }
    actions "delecs" {

    }
    actions "topics" {

    }
    actions "edit" {

    }
    actions "stmt" {

    }
  }