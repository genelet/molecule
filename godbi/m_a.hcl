tableName = "m_a"
pks     = ["id"]
idAuto  = "id"
uniques = ["x", "y"]
columns x {
  typeName    = "string"
  columnLabel = "x"
  notnull     = true
}

columns y {
  typeName    = "string"
  columnLabel = "y"
  notnull     = true
}

columns z {
  typeName    = "string"
  columnLabel = "z"
}

columns id {
  typeName    = "int"
  columnLabel = "id"
  auto        = true
}

actions topics {
  nextpages m_a edit {
    relateExtra = {
      id = "id"
    }
  }
  
}

actions insert {
  nextpages m_b insert {
    relateArgs = {
      id = "id"
    }
  }
  
}

actions insupd {
  nextpages m_b insert {
    relateArgs = {
      id = "id"
    }
  }
  
}

actions edit {
  nextpages m_b topics {
    relateExtra = {
      id = "id"
    }
  }
  
}

