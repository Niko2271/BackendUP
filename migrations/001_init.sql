    _, err = db.Exec(`
    create table if not exists roles(
        id text primary key,
        name text not null unique
    );`)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    insert or ignore into roles (id, name) values 
        ('1', 'admin'),
        ('2', 'manager'),
        ('3', 'client');
    `)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    create table if not exists users(
        id text primary key,
        name text not null,
        email text not null unique,
        phone text,
        roleId text not null,
        registrationDate datetime default current_timestamp,
        foreign key(roleId) references roles(id)
    );`)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    create table if not exists products(
        id text primary key,
        name text not null,
        type text,
        gender text,
        size text,
        color text,
        description text,
        brand text,
        season text,
        price real not null 
    );`)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    create table if not exists carts(
        id text primary key,
        userId text not null,
        productId text not null,
        count integer not null default 1,
        foreign key(userId) references users(id),
        foreign key(productId) references products(id)
    )`)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    create table if not exists orders(
        id text primary key,
        userId text not null,
        datetime text not null,
        status text not null,
        foreign key(userId) references users(id)
    )`)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    create table if not exists order_items(
        id text primary key,
        orderId text not null,
        productId text not null,
        price real not null,
        foreign key(orderId) references orders(id),
        foreign key(productId) references products(id)
    );`)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    create table if not exists photos(
        id text primary key,
        filename text not null,
        productId text not null,
        foreign key(productId) references products(id)
    )`)
