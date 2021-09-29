/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
  '/dataproducts': {
    /** List all dataproducts */
    get: {
      responses: {
        /** OK */
        200: {
          content: {
            'application/json': components['schemas']['Dataproduct'][]
          }
        }
      }
    }
    /** Create a new dataproduct */
    post: {
      responses: {
        /** Created successfully */
        201: {
          content: {
            'application/json': components['schemas']['Dataproduct']
          }
        }
      }
      requestBody: {
        content: {
          'application/json': components['schemas']['NewDataproduct']
        }
      }
    }
  }
  '/dataproducts/{dataproduct_id}': {
    /** List a dataproduct with datasets */
    get: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
        }
      }
      responses: {
        /** OK */
        200: {
          content: {
            'application/json': components['schemas']['Dataproduct'][]
          }
        }
      }
    }
    /** Update a dataproduct */
    put: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
        }
      }
      responses: {
        /** Updated OK */
        200: {
          content: {
            'application/json': components['schemas']['Dataproduct']
          }
        }
      }
      requestBody: {
        content: {
          'application/json': components['schemas']['NewDataproduct']
        }
      }
    }
    /** Delete a dataproduct */
    delete: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
        }
      }
      responses: {
        /** Deleted OK */
        204: never
      }
    }
  }
  '/dataproducts/{dataproduct_id}/datasets': {
    /** List all datasets for a dataproduct */
    get: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
        }
      }
      responses: {
        /** OK */
        200: {
          content: {
            'application/json': components['schemas']['Dataset'][]
          }
        }
      }
    }
    /** Create a new dataset */
    post: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
        }
      }
      responses: {
        /** Created successfully */
        201: {
          content: {
            'application/json': components['schemas']['Dataset']
          }
        }
      }
      requestBody: {
        content: {
          'application/json': components['schemas']['NewDataset']
        }
      }
    }
  }
  '/dataproducts/{dataproduct_id}/datasets/{dataset_id}': {
    /** Get dataset */
    get: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
          /** Dataset ID */
          dataset_id: string
        }
      }
      responses: {
        /** OK */
        200: {
          content: {
            'application/json': components['schemas']['Dataset']
          }
        }
      }
    }
    /** Update a dataset */
    put: {
      parameters: {
        path: {
          /** Dataproduct ID */
          dataproduct_id: string
          /** Dataset ID */
          dataset_id: string
        }
      }
      responses: {
        /** Updated OK */
        200: {
          content: {
            'application/json': components['schemas']['Dataset']
          }
        }
      }
      requestBody: {
        content: {
          'application/json': components['schemas']['NewDataset']
        }
      }
    }
    /** Delete a dataset */
    delete: operations['deleteDataset']
  }
  '/search': {
    get: {
      parameters: {
        query: {
          q?: string
        }
      }
      responses: {
        /** Search result */
        200: {
          content: {
            'application/json': components['schemas']['SearchResultEntry'][]
          }
        }
      }
    }
    parameters: {
      query: {
        q?: string
      }
    }
  }
}

export interface components {
  schemas: {
    Dataproduct: {
      id: string
      name: string
      description?: string
      slug: string
      repo?: string
      last_modified: string
      created: string
      owner: components['schemas']['Owner']
      keyword?: string[]
      datasets?: {
        id?: string
        type?: components['schemas']['DatasetType']
      }[]
    }
    DatasetType: 'bigquery'
    Owner: {
      team: string
      teamkatalogen?: string
    }
    NewDataproduct: {
      name: string
      description?: string
      slug?: string
      repo?: string
      owner: components['schemas']['Owner']
      keyword?: string[]
    }
    Dataset: {
      id?: string
      dataproduct_id?: string
      name?: string
      description?: string
      pii?: boolean
      bigquery?: {
        project_id: string
        dataset: string
        table: string
      }
    }
    NewDataset: {
      name: string
      description?: string
      pii: boolean
      bigquery: {
        project_id: string
        dataset: string
        table: string
      }
    }
    SearchResultEntry: {
      url?: string
      type?: 'dataset' | 'dataproduct' | 'datapackage'
      id?: string
      name?: string
      excerpt?: string
    }
  }
}

export interface operations {
  /** Delete a dataset */
  deleteDataset: {
    parameters: {
      path: {
        /** Dataproduct ID */
        dataproduct_id: string
        /** Dataset ID */
        dataset_id: string
      }
    }
    responses: {
      /** Deleted OK */
      204: never
    }
  }
}

export interface external {}
