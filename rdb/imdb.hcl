dbDriver = 2
atoms actor {
  tableName = "actor"
  pks       = ["aid"]
  idAuto    = "rowid"
  columns aid {
    typeName    = "int"
    columnLabel = "aid"
  }
  
  columns gender {
    typeName    = "string"
    columnLabel = "gender"
  }
  
  columns name {
    typeName    = "string"
    columnLabel = "name"
  }
  
  columns nationality {
    typeName    = "string"
    columnLabel = "nationality"
  }
  
  columns birth_city {
    typeName    = "string"
    columnLabel = "birth_city"
  }
  
  columns birth_year {
    typeName    = "int"
    columnLabel = "birth_year"
  }
  
  actions edit {
    nextpages cast topics {
      relateExtra = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions topics {
    nextpages cast topics {
      relateExtra = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions read {
    nextpages cast list {
      relateExtra = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions list {
    nextpages cast list {
      relateExtra = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions insert {
    nextpages cast insert {
      relateArgs = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions update {
    nextpages cast update {
      relateArgs = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions insupd {
    nextpages cast insupd {
      relateArgs = {
        aid = "aid"
      }
      marker = "cast"
    }
    
  }
  
  actions delete {
    prepares cast delecs {
      relateArgs = {
        aid = "aid"
      }
    }
    
  }
  
  actions delecs {
    nextpages actor delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms copyright {
  tableName = "copyright"
  pks       = ["id"]
  idAuto    = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns cid {
    typeName    = "int"
    columnLabel = "cid"
  }
  
  actions edit {
    nextpages cast topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions topics {
    nextpages cast topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by topics {
      relateExtra = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions read {
    nextpages cast list {
      relateExtra = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification list {
      relateExtra = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by list {
      relateExtra = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by list {
      relateExtra = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags list {
      relateExtra = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by list {
      relateExtra = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions list {
    nextpages cast list {
      relateExtra = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification list {
      relateExtra = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by list {
      relateExtra = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by list {
      relateExtra = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags list {
      relateExtra = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by list {
      relateExtra = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions insert {
    nextpages cast insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions update {
    nextpages cast update {
      relateArgs = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification update {
      relateArgs = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by update {
      relateArgs = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by update {
      relateArgs = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags update {
      relateArgs = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by update {
      relateArgs = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions insupd {
    nextpages cast insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "cast"
    }
    
    nextpages classification insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "classification"
    }
    
    nextpages directed_by insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "directed_by"
    }
    
    nextpages made_by insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "made_by"
    }
    
    nextpages tags insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "tags"
    }
    
    nextpages written_by insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "written_by"
    }
    
  }
  
  actions delete {
    prepares cast delecs {
      relateArgs = {
        msid = "msid"
      }
    }
    
    prepares classification delecs {
      relateArgs = {
        msid = "msid"
      }
    }
    
    prepares directed_by delecs {
      relateArgs = {
        msid = "msid"
      }
    }
    
    prepares made_by delecs {
      relateArgs = {
        msid = "msid"
      }
    }
    
    prepares tags delecs {
      relateArgs = {
        msid = "msid"
      }
    }
    
    prepares written_by delecs {
      relateArgs = {
        msid = "msid"
      }
    }
    
  }
  
  actions delecs {
    nextpages copyright delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms cast {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns aid {
    typeName    = "int"
    columnLabel = "aid"
  }
  
  columns role {
    typeName    = "int"
    columnLabel = "role"
  }
  
  fks {
    fkTable  = "copyright"
    fkColumn = "msid"
    column   = "msid"
  }
  
  fks {
    fkTable  = "actor"
    fkColumn = "aid"
    column   = "aid"
  }
  
  actions insert {
    prepares copyright insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
    prepares actor insert {
      relateArgs = {
        aid = "aid"
      }
      marker = "actor"
    }
    
  }
  
  actions update {
    prepares copyright update {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
    prepares actor update {
      relateArgs = {
        aid = "aid"
      }
      marker = "actor"
    }
    
  }
  
  actions insupd {
    prepares copyright insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
    prepares actor insupd {
      relateArgs = {
        aid = "aid"
      }
      marker = "actor"
    }
    
  }
  
  actions delecs {
    nextpages cast delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms genre {
  pks    = ["gid"]
  idAuto = "rowid"
  columns gid {
    typeName    = "int"
    columnLabel = "gid"
  }
  
  columns genre {
    typeName    = "string"
    columnLabel = "genre"
  }
  
  actions edit {
    nextpages classification topics {
      relateExtra = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions topics {
    nextpages classification topics {
      relateExtra = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions read {
    nextpages classification list {
      relateExtra = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions list {
    nextpages classification list {
      relateExtra = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions insert {
    nextpages classification insert {
      relateArgs = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions update {
    nextpages classification update {
      relateArgs = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions insupd {
    nextpages classification insupd {
      relateArgs = {
        gid = "gid"
      }
      marker = "classification"
    }
    
  }
  
  actions delete {
    prepares classification delecs {
      relateArgs = {
        gid = "gid"
      }
    }
    
  }
  
  actions delecs {
    nextpages genre delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms classification {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns gid {
    typeName    = "int"
    columnLabel = "gid"
  }
  
  fks {
    fkTable  = "copyright"
    fkColumn = "msid"
    column   = "msid"
  }
  
  fks {
    fkTable  = "genre"
    fkColumn = "gid"
    column   = "gid"
  }
  
  actions insert {
    prepares copyright insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
    prepares genre insert {
      relateArgs = {
        gid = "gid"
      }
      marker = "genre"
    }
    
  }
  
  actions update {
    prepares copyright update {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
    prepares genre update {
      relateArgs = {
        gid = "gid"
      }
      marker = "genre"
    }
    
  }
  
  actions insupd {
    prepares copyright insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
    prepares genre insupd {
      relateArgs = {
        gid = "gid"
      }
      marker = "genre"
    }
    
  }
  
  actions delecs {
    nextpages classification delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms company {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns name {
    typeName    = "string"
    columnLabel = "name"
  }
  
  columns country_code {
    typeName    = "string"
    columnLabel = "country_code"
  }
  
  actions delecs {
    nextpages company delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms director {
  pks    = ["did"]
  idAuto = "rowid"
  columns did {
    typeName    = "int"
    columnLabel = "did"
  }
  
  columns gender {
    typeName    = "string"
    columnLabel = "gender"
  }
  
  columns name {
    typeName    = "string"
    columnLabel = "name"
  }
  
  columns nationality {
    typeName    = "string"
    columnLabel = "nationality"
  }
  
  columns birth_city {
    typeName    = "string"
    columnLabel = "birth_city"
  }
  
  columns birth_year {
    typeName    = "int"
    columnLabel = "birth_year"
  }
  
  actions edit {
    nextpages directed_by topics {
      relateExtra = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions topics {
    nextpages directed_by topics {
      relateExtra = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions read {
    nextpages directed_by list {
      relateExtra = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions list {
    nextpages directed_by list {
      relateExtra = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions insert {
    nextpages directed_by insert {
      relateArgs = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions update {
    nextpages directed_by update {
      relateArgs = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions insupd {
    nextpages directed_by insupd {
      relateArgs = {
        did = "did"
      }
      marker = "directed_by"
    }
    
  }
  
  actions delete {
    prepares directed_by delecs {
      relateArgs = {
        did = "did"
      }
    }
    
  }
  
  actions delecs {
    nextpages director delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms producer {
  pks    = ["pid"]
  idAuto = "rowid"
  columns pid {
    typeName    = "int"
    columnLabel = "pid"
  }
  
  columns gender {
    typeName    = "string"
    columnLabel = "gender"
  }
  
  columns name {
    typeName    = "string"
    columnLabel = "name"
  }
  
  columns nationality {
    typeName    = "string"
    columnLabel = "nationality"
  }
  
  columns birth_city {
    typeName    = "string"
    columnLabel = "birth_city"
  }
  
  columns birth_year {
    typeName    = "int"
    columnLabel = "birth_year"
  }
  
  actions edit {
    nextpages made_by topics {
      relateExtra = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions topics {
    nextpages made_by topics {
      relateExtra = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions read {
    nextpages made_by list {
      relateExtra = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions list {
    nextpages made_by list {
      relateExtra = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions insert {
    nextpages made_by insert {
      relateArgs = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions update {
    nextpages made_by update {
      relateArgs = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions insupd {
    nextpages made_by insupd {
      relateArgs = {
        pid = "pid"
      }
      marker = "made_by"
    }
    
  }
  
  actions delete {
    prepares made_by delecs {
      relateArgs = {
        pid = "pid"
      }
    }
    
  }
  
  actions delecs {
    nextpages producer delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms directed_by {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns did {
    typeName    = "int"
    columnLabel = "did"
  }
  
  fks {
    fkTable  = "director"
    fkColumn = "did"
    column   = "did"
  }
  
  fks {
    fkTable  = "copyright"
    fkColumn = "msid"
    column   = "msid"
  }
  
  actions insert {
    prepares director insert {
      relateArgs = {
        did = "did"
      }
      marker = "director"
    }
    
    prepares copyright insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions update {
    prepares director update {
      relateArgs = {
        did = "did"
      }
      marker = "director"
    }
    
    prepares copyright update {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions insupd {
    prepares director insupd {
      relateArgs = {
        did = "did"
      }
      marker = "director"
    }
    
    prepares copyright insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions delecs {
    nextpages directed_by delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms keyword {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns keyword {
    typeName    = "string"
    columnLabel = "keyword"
  }
  
  actions edit {
    nextpages tags topics {
      relateExtra = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions topics {
    nextpages tags topics {
      relateExtra = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions read {
    nextpages tags list {
      relateExtra = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions list {
    nextpages tags list {
      relateExtra = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions insert {
    nextpages tags insert {
      relateArgs = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions update {
    nextpages tags update {
      relateArgs = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions insupd {
    nextpages tags insupd {
      relateArgs = {
        kid = "kid"
      }
      marker = "tags"
    }
    
  }
  
  actions delete {
    prepares tags delecs {
      relateArgs = {
        kid = "kid"
      }
    }
    
  }
  
  actions delecs {
    nextpages keyword delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms made_by {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns pid {
    typeName    = "int"
    columnLabel = "pid"
  }
  
  fks {
    fkTable  = "producer"
    fkColumn = "pid"
    column   = "pid"
  }
  
  fks {
    fkTable  = "copyright"
    fkColumn = "msid"
    column   = "msid"
  }
  
  actions insert {
    prepares producer insert {
      relateArgs = {
        pid = "pid"
      }
      marker = "producer"
    }
    
    prepares copyright insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions update {
    prepares producer update {
      relateArgs = {
        pid = "pid"
      }
      marker = "producer"
    }
    
    prepares copyright update {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions insupd {
    prepares producer insupd {
      relateArgs = {
        pid = "pid"
      }
      marker = "producer"
    }
    
    prepares copyright insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions delecs {
    nextpages made_by delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms movie {
  pks    = ["mid"]
  idAuto = "rowid"
  columns mid {
    typeName    = "int"
    columnLabel = "mid"
  }
  
  columns title {
    typeName    = "string"
    columnLabel = "title"
  }
  
  columns release_year {
    typeName    = "int"
    columnLabel = "release_year"
  }
  
  columns title_aka {
    typeName    = "string"
    columnLabel = "title_aka"
  }
  
  columns budget {
    typeName    = "string"
    columnLabel = "budget"
  }
  
  actions delecs {
    nextpages movie delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms tags {
  pks    = ["id"]
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns kid {
    typeName    = "int"
    columnLabel = "kid"
  }
  
  fks {
    fkTable  = "keyword"
    fkColumn = "kid"
    column   = "kid"
  }
  
  fks {
    fkTable  = "copyright"
    fkColumn = "msid"
    column   = "msid"
  }
  
  actions insert {
    prepares keyword insert {
      relateArgs = {
        kid = "kid"
      }
      marker = "keyword"
    }
    
    prepares copyright insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions update {
    prepares keyword update {
      relateArgs = {
        kid = "kid"
      }
      marker = "keyword"
    }
    
    prepares copyright update {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions insupd {
    prepares keyword insupd {
      relateArgs = {
        kid = "kid"
      }
      marker = "keyword"
    }
    
    prepares copyright insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions delecs {
    nextpages tags delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms tv_series {
  pks    = ["sid"]
  idAuto = "rowid"
  columns sid {
    typeName    = "int"
    columnLabel = "sid"
  }
  
  columns title {
    typeName    = "string"
    columnLabel = "title"
  }
  
  columns release_year {
    typeName    = "int"
    columnLabel = "release_year"
  }
  
  columns num_of_seasons {
    typeName    = "int"
    columnLabel = "num_of_seasons"
  }
  
  columns num_of_episodes {
    typeName    = "int"
    columnLabel = "num_of_episodes"
  }
  
  columns title_aka {
    typeName    = "string"
    columnLabel = "title_aka"
  }
  
  columns budget {
    typeName    = "string"
    columnLabel = "budget"
  }
  
  actions delecs {
    nextpages tv_series delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms writer {
  pks    = ["wid"]
  idAuto = "rowid"
  columns wid {
    typeName    = "int"
    columnLabel = "wid"
  }
  
  columns gender {
    typeName    = "string"
    columnLabel = "gender"
  }
  
  columns name {
    typeName    = "int"
    columnLabel = "name"
  }
  
  columns nationality {
    typeName    = "int"
    columnLabel = "nationality"
  }
  
  columns num_of_episodes {
    typeName    = "int"
    columnLabel = "num_of_episodes"
  }
  
  columns birth_city {
    typeName    = "string"
    columnLabel = "birth_city"
  }
  
  columns birth_year {
    typeName    = "int"
    columnLabel = "birth_year"
  }
  
  actions edit {
    nextpages written_by topics {
      relateExtra = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions topics {
    nextpages written_by topics {
      relateExtra = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions read {
    nextpages written_by list {
      relateExtra = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions list {
    nextpages written_by list {
      relateExtra = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions insert {
    nextpages written_by insert {
      relateArgs = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions update {
    nextpages written_by update {
      relateArgs = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions insupd {
    nextpages written_by insupd {
      relateArgs = {
        wid = "wid"
      }
      marker = "written_by"
    }
    
  }
  
  actions delete {
    prepares written_by delecs {
      relateArgs = {
        wid = "wid"
      }
    }
    
  }
  
  actions delecs {
    nextpages writer delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

atoms written_by {
  idAuto = "rowid"
  columns id {
    typeName    = "int"
    columnLabel = "id"
  }
  
  columns msid {
    typeName    = "int"
    columnLabel = "msid"
  }
  
  columns wid {
    typeName    = "int"
    columnLabel = "wid"
  }
  
  fks {
    fkTable  = "writer"
    fkColumn = "wid"
    column   = "wid"
  }
  
  fks {
    fkTable  = "copyright"
    fkColumn = "msid"
    column   = "msid"
  }
  
  actions insert {
    prepares writer insert {
      relateArgs = {
        wid = "wid"
      }
      marker = "writer"
    }
    
    prepares copyright insert {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions update {
    prepares writer update {
      relateArgs = {
        wid = "wid"
      }
      marker = "writer"
    }
    
    prepares copyright update {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions insupd {
    prepares writer insupd {
      relateArgs = {
        wid = "wid"
      }
      marker = "writer"
    }
    
    prepares copyright insupd {
      relateArgs = {
        msid = "msid"
      }
      marker = "copyright"
    }
    
  }
  
  actions delecs {
    nextpages written_by delete {
      relateArgs = {
        rowid = "rowid"
      }
    }
    
  }
  
}

