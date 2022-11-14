CREATE PROCEDURE RegistrarMetodoPago



@CodMetodoPago VARCHAR(7),


@Descripcion VARCHAR(20),


@Premiacion BIT


AS


BEGIN


	INSERT INTO MetodosPago(CodMetodoPago, Descripcion, Premiacion) VALUES(@CodMetodoPago, @Descripcion, @Premiacion)


END

GO
