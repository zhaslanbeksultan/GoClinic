# GoClinic Project  
## Restaurants REST API  
POST /menus  
GET /menus/:id  
PUT /menus/:id  
DELETE /menus/:id  

## DB Structure  
Table doctors {  
  id bigserial [primary key]  
  created_at timestamp  
  updated_at timestamp  
  first_name text  
  last_name text  
  speciality text  
  phone text  
} 

Table patients {  
  id bigserial [primary key]  
  created_at timestamp  
  updated_at timestamp  
  first_name text  
  last_name text  
  phone text  
}  

// many-to-many  
Table appointments {  
  id bigserial [primary key]  
  created_at timestamp  
  updated_at timestamp  
  doctor_id bigserial  
  patient_id bigserial  
}  

Ref: appointments.doctor_id < doctors.id  
Ref: appointments.patient_id < patients.id  
